package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hanayo-dot/KIVU/backend/config"
	"github.com/hanayo-dot/KIVU/backend/integrations"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// CopernicusSyncService periodically synchronizes satellite-derived LWQ and LSWT metrics.
type CopernicusSyncService struct {
	cfg    *config.Config
	db     *sqlx.DB
	client *integrations.CopernicusClient
}

// NewCopernicusSyncService initializes CopernicusSyncService.
func NewCopernicusSyncService(cfg *config.Config, db *sqlx.DB, client *integrations.CopernicusClient) *CopernicusSyncService {
	return &CopernicusSyncService{
		cfg:    cfg,
		db:     db,
		client: client,
	}
}

// StartSyncLoop starts the 24-hour background ticker sync.
func (s *CopernicusSyncService) StartSyncLoop() {
	if !s.cfg.UseCopernicusLive {
		log.Println("[Copernicus Sync] Live satellite sync is disabled (USE_COPERNICUS_LIVE=false). Operating in seed/fallback mode.")
		return
	}

	log.Println("[Copernicus Sync] Starting 24h satellite telemetry synchronization loop...")
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			s.SyncZoneMetrics()
			<-ticker.C
		}
	}()
}

// SyncZoneMetrics fetches Sentinel Hub stats and updates zone_metrics.
func (s *CopernicusSyncService) SyncZoneMetrics() {
	var zones []models.LakeZone
	err := s.db.Select(&zones, "SELECT id, name, boundary, region_label FROM lake_zones")
	if err != nil {
		log.Printf("[Copernicus Sync] Failed to query lake zones: %v", err)
		return
	}

	for _, zone := range zones {
		var polygon models.GeoJSONPolygon
		if err := json.Unmarshal(zone.Boundary, &polygon); err != nil {
			log.Printf("[Copernicus Sync] Invalid GeoJSON for zone %s: %v", zone.ID, err)
			continue
		}

		metrics, err := s.client.FetchZoneStatistics(polygon)
		if err != nil {
			log.Printf("[Copernicus Sync] Failed to fetch stats for zone %s: %v", zone.ID, err)
			continue
		}

		// Save the metrics to the database
		_, err = s.db.Exec(`
			INSERT INTO zone_metrics (id, zone_id, period, avg_dissolved_oxygen, avg_temperature, avg_ph, avg_turbidity, risk_level, trend)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			uuid.New(), zone.ID, metrics.Period, metrics.AvgDissolvedOxygen, metrics.AvgTemperature, metrics.AvgPH, metrics.AvgTurbidity, metrics.RiskLevel, metrics.Trend)

		if err != nil {
			log.Printf("[Copernicus Sync] Failed to insert metrics for zone %s: %v", zone.ID, err)
			continue
		}
		log.Printf("[Copernicus Sync] Successfully synced metrics for zone %s", zone.ID)
	}

	log.Println("[Copernicus Sync] Completed satellite telemetry synchronization loop.")
}
