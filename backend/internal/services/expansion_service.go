package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// ExpansionService evaluates lake zones for aquaculture expansion suitability.
type ExpansionService struct {
	db *sqlx.DB
}

// NewExpansionService initializes ExpansionService.
func NewExpansionService(db *sqlx.DB) *ExpansionService {
	return &ExpansionService{db: db}
}

// GetAllSignals fetches current expansion suitability signals across all zones.
func (s *ExpansionService) GetAllSignals() ([]models.ExpansionSignal, error) {
	var signals []models.ExpansionSignal
	query := `
		SELECT DISTINCT ON (zone_id) id, zone_id, suitability, rationale, computed_at 
		FROM expansion_signals 
		ORDER BY zone_id, computed_at DESC`

	err := s.db.Select(&signals, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expansion signals: %w", err)
	}
	return signals, nil
}

// GetSignalByZone fetches the latest suitability signal for a given zone.
func (s *ExpansionService) GetSignalByZone(zoneID uuid.UUID) (*models.ExpansionSignal, error) {
	var signal models.ExpansionSignal
	query := `
		SELECT id, zone_id, suitability, rationale, computed_at 
		FROM expansion_signals 
		WHERE zone_id = $1 
		ORDER BY computed_at DESC LIMIT 1`

	err := s.db.Get(&signal, query, zoneID)
	if err != nil {
		return nil, fmt.Errorf("no expansion signal found for zone %s: %w", zoneID, err)
	}
	return &signal, nil
}

// EvaluateCoordinates takes a lat/lon and uses PostGIS ST_Distance to return the nearest zone's suitability.
func (s *ExpansionService) EvaluateCoordinates(lat, lon float64) (*models.ExpansionSignal, error) {
	var zoneID uuid.UUID
	spatialQuery := `
		SELECT id 
		FROM lake_zones 
		ORDER BY ST_Distance(boundary, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) ASC 
		LIMIT 1`

	err := s.db.Get(&zoneID, spatialQuery, lon, lat)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearest lake zone for coordinates (%.4f, %.4f): %w", lat, lon, err)
	}

	return s.GetSignalByZone(zoneID)
}
