package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// LakeHandler manages Level 3 Lake Victoria spatial intelligence endpoints.
type LakeHandler struct {
	db *sqlx.DB
}

// NewLakeHandler initializes LakeHandler.
func NewLakeHandler(db *sqlx.DB) *LakeHandler {
	return &LakeHandler{db: db}
}

// GetLakeZones handles GET /lake/zones. Returns PostGIS ST_AsGeoJSON boundaries.
func (h *LakeHandler) GetLakeZones(c *gin.Context) {
	type ZoneWithMetrics struct {
		models.LakeZone
		CurrentMetrics *models.ZoneMetrics `json:"current_metrics"`
	}

	query := `
		SELECT 
			z.id, z.name, z.region_label,
			ST_AsGeoJSON(z.boundary)::json AS boundary,
			zm.id AS "current_metrics.id",
			zm.zone_id AS "current_metrics.zone_id",
			COALESCE(zm.period, '') AS "current_metrics.period",
			COALESCE(zm.avg_dissolved_oxygen, 0.0) AS "current_metrics.avg_dissolved_oxygen",
			COALESCE(zm.avg_temperature, 0.0) AS "current_metrics.avg_temperature",
			COALESCE(zm.avg_ph, 0.0) AS "current_metrics.avg_ph",
			COALESCE(zm.avg_turbidity, 0.0) AS "current_metrics.avg_turbidity",
			COALESCE(zm.risk_level, 'low') AS "current_metrics.risk_level",
			COALESCE(zm.trend, 'stable') AS "current_metrics.trend",
			COALESCE(zm.computed_at, NOW()) AS "current_metrics.computed_at"
		FROM lake_zones z
		LEFT JOIN LATERAL (
			SELECT * FROM zone_metrics WHERE zone_id = z.id ORDER BY computed_at DESC LIMIT 1
		) zm ON TRUE`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to query lake zones: "+err.Error(), 500))
		return
	}
	defer rows.Close()

	var results []ZoneWithMetrics
	for rows.Next() {
		var item ZoneWithMetrics
		var boundaryRaw []byte
		var zm models.ZoneMetrics

		err := rows.Scan(
			&item.ID, &item.Name, &item.RegionLabel, &boundaryRaw,
			&zm.ID, &zm.ZoneID, &zm.Period, &zm.AvgDissolvedOxygen,
			&zm.AvgTemperature, &zm.AvgPH, &zm.AvgTurbidity,
			&zm.RiskLevel, &zm.Trend, &zm.ComputedAt,
		)
		if err != nil {
			continue
		}
		item.Boundary = json.RawMessage(boundaryRaw)
		if zm.Period != "" {
			item.CurrentMetrics = &zm
		}
		results = append(results, item)
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(results))
}

// GetZoneByID handles GET /lake/zones/:id.
func (h *LakeHandler) GetZoneByID(c *gin.Context) {
	zoneID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid zone ID parameter", 400))
		return
	}

	var zone models.LakeZone
	var boundaryRaw []byte
	err = h.db.QueryRow(`SELECT id, name, region_label, ST_AsGeoJSON(boundary)::json FROM lake_zones WHERE id = $1`, zoneID).
		Scan(&zone.ID, &zone.Name, &zone.RegionLabel, &boundaryRaw)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse("Zone not found", 404))
		return
	}
	zone.Boundary = json.RawMessage(boundaryRaw)

	var metrics models.ZoneMetrics
	var metricsPtr *models.ZoneMetrics
	errM := h.db.Get(&metrics, `SELECT id, zone_id, period, avg_dissolved_oxygen, avg_temperature, avg_ph, avg_turbidity, risk_level, trend, computed_at FROM zone_metrics WHERE zone_id = $1 ORDER BY computed_at DESC LIMIT 1`, zoneID)
	if errM == nil {
		metricsPtr = &metrics
	}

	var signal models.ExpansionSignal
	var signalPtr *models.ExpansionSignal
	errS := h.db.Get(&signal, `SELECT id, zone_id, suitability, rationale, computed_at FROM expansion_signals WHERE zone_id = $1 ORDER BY computed_at DESC LIMIT 1`, zoneID)
	if errS == nil {
		signalPtr = &signal
	}

	var activeCages, incCount int
	_ = h.db.Get(&activeCages, `SELECT COUNT(*) FROM cages WHERE ST_Contains((SELECT boundary::geometry FROM lake_zones WHERE id = $1), location::geometry)`, zoneID)
	_ = h.db.Get(&incCount, `SELECT COUNT(*) FROM incidents i JOIN cages c ON i.cage_id = c.id WHERE ST_Contains((SELECT boundary::geometry FROM lake_zones WHERE id = $1), c.location::geometry)`, zoneID)

	detail := models.LakeZoneDetail{
		Zone:            zone,
		CurrentMetrics:  metricsPtr,
		ExpansionSignal: signalPtr,
		IncidentCount:   incCount,
		ActiveCageCount: activeCages,
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(detail))
}

// GetRegionalAlerts handles GET /lake/alerts.
func (h *LakeHandler) GetRegionalAlerts(c *gin.Context) {
	var alerts []models.Alert
	err := h.db.Select(&alerts, `SELECT id, scope, related_id, severity, message, triggered_at, acknowledged FROM alerts WHERE scope = 'region' ORDER BY triggered_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch regional alerts: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(alerts))
}

// GetHotspots handles GET /lake/hotspots. Returns zones flagged with risk_level = 'high'.
func (h *LakeHandler) GetHotspots(c *gin.Context) {
	type HotspotZone struct {
		ZoneID             uuid.UUID `json:"zone_id"`
		Name               string    `json:"name"`
		RegionLabel        string    `json:"region_label"`
		AvgDissolvedOxygen float64   `json:"avg_dissolved_oxygen"`
		AvgTemperature     float64   `json:"avg_temperature"`
		RiskLevel          string    `json:"risk_level"`
		Trend              string    `json:"trend"`
	}

	query := `
		SELECT z.id AS zone_id, z.name, z.region_label, zm.avg_dissolved_oxygen, zm.avg_temperature, zm.risk_level, zm.trend
		FROM lake_zones z
		JOIN zone_metrics zm ON z.id = zm.zone_id
		WHERE zm.risk_level = 'high'
		ORDER BY zm.computed_at DESC`

	var hotspots []HotspotZone
	err := h.db.Select(&hotspots, query)
	if err != nil {
		hotspots = []HotspotZone{}
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(hotspots))
}
