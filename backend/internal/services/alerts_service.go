package services

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/hanayo-dot/KIVU/backend/models"
)

// AlertsService handles multi-tier alert triggers and regional escalation.
type AlertsService struct {
	db *sqlx.DB
}

// NewAlertsService initializes AlertsService.
func NewAlertsService(db *sqlx.DB) *AlertsService {
	return &AlertsService{db: db}
}

// TriggerCageAlert creates a cage-scoped alert and checks for farm/region escalation.
func (s *AlertsService) TriggerCageAlert(cageID uuid.UUID, severity, message string) error {
	// 1. Insert cage-scoped alert
	query := `INSERT INTO alerts (scope, related_id, severity, message) VALUES ('cage', $1, $2, $3)`
	_, err := s.db.Exec(query, cageID.String(), severity, message)
	if err != nil {
		return fmt.Errorf("failed to insert cage alert: %w", err)
	}

	// 2. Check Farm-level escalation
	var farmID uuid.UUID
	err = s.db.Get(&farmID, `SELECT farm_id FROM cages WHERE id = $1`, cageID)
	if err == nil {
		var activeCageAlertCount int
		_ = s.db.Get(&activeCageAlertCount,
			`SELECT COUNT(DISTINCT related_id) 
			 FROM alerts 
			 WHERE scope = 'cage' 
			   AND triggered_at >= NOW() - INTERVAL '24 hours' 
			   AND related_id::uuid IN (SELECT id FROM cages WHERE farm_id = $1)`,
			farmID)

		if activeCageAlertCount >= 2 {
			farmMsg := fmt.Sprintf("FARM ESCALATION: Multiple cages (%d) in farm %s reporting environmental anomalies", activeCageAlertCount, farmID)
			_, _ = s.db.Exec(`INSERT INTO alerts (scope, related_id, severity, message) VALUES ('farm', $1, 'warning', $2)`, farmID.String(), farmMsg)
			log.Printf("[Alert Escalation] Escalated to Farm %s scope", farmID)
		}
	}

	// 3. Check Regional escalation across lake zones
	var zoneName string
	zoneQuery := `
		SELECT z.name 
		FROM lake_zones z 
		JOIN cages c ON ST_Contains(z.boundary::geometry, c.location::geometry) 
		WHERE c.id = $1 LIMIT 1`
	errZone := s.db.Get(&zoneName, zoneQuery, cageID)
	if errZone == nil && zoneName != "" {
		var activeFarmAlerts int
		_ = s.db.Get(&activeFarmAlerts,
			`SELECT COUNT(DISTINCT a.related_id) 
			 FROM alerts a 
			 JOIN cages c ON a.related_id::uuid = c.id 
			 JOIN lake_zones z ON ST_Contains(z.boundary::geometry, c.location::geometry) 
			 WHERE a.scope = 'cage' 
			   AND a.triggered_at >= NOW() - INTERVAL '24 hours' 
			   AND z.name = $1`,
			zoneName)

		if activeFarmAlerts >= 2 {
			regionMsg := fmt.Sprintf("REGIONAL ALERT: Widespread environmental risk detected across multiple farms in Lake Zone '%s'", zoneName)
			_, _ = s.db.Exec(`INSERT INTO alerts (scope, related_id, severity, message) VALUES ('region', $1, 'critical', $2)`, zoneName, regionMsg)
			log.Printf("[Alert Escalation] Escalated to Region '%s' scope", zoneName)
		}
	}

	return nil
}

// GetAlerts fetches alerts filtered optional scope parameter.
func (s *AlertsService) GetAlerts(scope string) ([]models.Alert, error) {
	var alerts []models.Alert
	var err error

	if scope != "" {
		err = s.db.Select(&alerts, `SELECT id, scope, related_id, severity, message, triggered_at, acknowledged FROM alerts WHERE scope = $1 ORDER BY triggered_at DESC`, scope)
	} else {
		err = s.db.Select(&alerts, `SELECT id, scope, related_id, severity, message, triggered_at, acknowledged FROM alerts ORDER BY triggered_at DESC`)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch alerts: %w", err)
	}
	return alerts, nil
}

// AcknowledgeAlert sets acknowledged = true.
func (s *AlertsService) AcknowledgeAlert(alertID uuid.UUID) error {
	res, err := s.db.Exec(`UPDATE alerts SET acknowledged = TRUE WHERE id = $1`, alertID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("alert not found")
	}
	return nil
}
