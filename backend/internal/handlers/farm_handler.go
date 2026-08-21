package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/hanayo-dot/KIVU/backend/internal/middleware"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// FarmHandler manages Level 2 My Farm endpoints.
type FarmHandler struct {
	db *sqlx.DB
}

// NewFarmHandler initializes FarmHandler.
func NewFarmHandler(db *sqlx.DB) *FarmHandler {
	return &FarmHandler{db: db}
}

func (h *FarmHandler) verifyFarmOwnership(c *gin.Context, farmID uuid.UUID) bool {
	farmerID, ok := middleware.GetAuthenticatedFarmerID(c)
	if !ok {
		return false
	}
	var count int
	_ = h.db.Get(&count, `SELECT COUNT(*) FROM farms WHERE id = $1 AND farmer_id = $2`, farmID, farmerID)
	return count > 0
}

// GetFarmDashboard handles GET /farms/:id/dashboard.
func (h *FarmHandler) GetFarmDashboard(c *gin.Context) {
	farmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid farm ID parameter", 400))
		return
	}

	if !h.verifyFarmOwnership(c, farmID) {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Access denied: Farm does not belong to farmer", 403))
		return
	}

	var farm models.Farm
	err = h.db.Get(&farm, `SELECT id, farmer_id, name, region, created_at FROM farms WHERE id = $1`, farmID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse("Farm not found", 404))
		return
	}

	var totalCages, activeCages, openIncidents int
	_ = h.db.Get(&totalCages, `SELECT COUNT(*) FROM cages WHERE farm_id = $1`, farmID)
	_ = h.db.Get(&activeCages, `SELECT COUNT(*) FROM cages WHERE farm_id = $1 AND status = 'active'`, farmID)
	_ = h.db.Get(&openIncidents, `SELECT COUNT(*) FROM incidents i JOIN cages c ON i.cage_id = c.id WHERE c.farm_id = $1 AND i.status = 'open'`, farmID)

	var recentAlerts []models.Alert
	_ = h.db.Select(&recentAlerts,
		`SELECT a.id, a.scope, a.related_id, a.severity, a.message, a.triggered_at, a.acknowledged 
		 FROM alerts a 
		 JOIN cages c ON a.related_id::uuid = c.id 
		 WHERE c.farm_id = $1 
		 ORDER BY a.triggered_at DESC LIMIT 5`, farmID)

	var cageSummaries []models.CageStatusSummary
	query := `
		SELECT 
			c.id AS cage_id, c.name, c.status,
			COALESCE(r.dissolved_oxygen, 0.0) AS latest_dissolved_oxygen,
			COALESCE(r.temperature, 0.0) AS latest_temperature,
			COALESCE(r.ph, 0.0) AS latest_ph,
			COALESCE(r.turbidity, 0.0) AS latest_turbidity,
			COALESCE(r.recorded_at, c.installed_at) AS last_recorded_at
		FROM cages c
		LEFT JOIN LATERAL (
			SELECT dissolved_oxygen, temperature, ph, turbidity, recorded_at
			FROM readings WHERE cage_id = c.id
			ORDER BY recorded_at DESC LIMIT 1
		) r ON TRUE
		WHERE c.farm_id = $1`

	err = h.db.Select(&cageSummaries, query, farmID)
	if err != nil {
		cageSummaries = []models.CageStatusSummary{}
	}

	view := models.FarmDashboardView{
		Farm:          farm,
		TotalCages:    totalCages,
		ActiveCages:   activeCages,
		OpenIncidents: openIncidents,
		RecentAlerts:  recentAlerts,
		CageStatuses:  cageSummaries,
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(view))
}

// CompareCages handles GET /farms/:id/cages/compare.
func (h *FarmHandler) CompareCages(c *gin.Context) {
	farmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid farm ID parameter", 400))
		return
	}

	if !h.verifyFarmOwnership(c, farmID) {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Access denied: Farm does not belong to farmer", 403))
		return
	}

	var cages []models.Cage
	err = h.db.Select(&cages, `SELECT id, farm_id, name, latitude, longitude, installed_at, status FROM cages WHERE farm_id = $1`, farmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch cages: "+err.Error(), 500))
		return
	}

	if len(cages) == 0 {
		c.JSON(http.StatusOK, models.NewSuccessResponse([]models.CageComparisonItem{}))
		return
	}

	// 1. Batch fetch latest readings per cage
	var readings []models.Reading
	_ = h.db.Select(&readings, `
		SELECT DISTINCT ON (cage_id) id, cage_id, dissolved_oxygen, temperature, ph, turbidity, recorded_at, source
		FROM readings
		WHERE cage_id IN (SELECT id FROM cages WHERE farm_id = $1)
		ORDER BY cage_id, recorded_at DESC`, farmID)

	readingMap := make(map[uuid.UUID]*models.Reading)
	for i := range readings {
		r := readings[i]
		readingMap[r.CageID] = &r
	}

	// 2. Batch fetch incident counts
	type countMap struct {
		CageID uuid.UUID `db:"cage_id"`
		Count  int       `db:"count"`
	}
	var incCounts []countMap
	_ = h.db.Select(&incCounts, `
		SELECT cage_id, COUNT(*) AS count
		FROM incidents
		WHERE cage_id IN (SELECT id FROM cages WHERE farm_id = $1)
		GROUP BY cage_id`, farmID)
	incMap := make(map[uuid.UUID]int)
	for _, ic := range incCounts {
		incMap[ic.CageID] = ic.Count
	}

	// 3. Batch fetch alert counts
	var alertCounts []countMap
	_ = h.db.Select(&alertCounts, `
		SELECT related_id::uuid AS cage_id, COUNT(*) AS count
		FROM alerts
		WHERE scope = 'cage' AND related_id::uuid IN (SELECT id FROM cages WHERE farm_id = $1)
		GROUP BY related_id::uuid`, farmID)
	alertMap := make(map[uuid.UUID]int)
	for _, ac := range alertCounts {
		alertMap[ac.CageID] = ac.Count
	}

	// 4. Batch fetch 7-day averages
	type avgMap struct {
		CageID  uuid.UUID `db:"cage_id"`
		DOAvg   float64   `db:"do_avg"`
		TempAvg float64   `db:"temp_avg"`
	}
	var avgs []avgMap
	_ = h.db.Select(&avgs, `
		SELECT cage_id, COALESCE(AVG(dissolved_oxygen), 0.0) AS do_avg, COALESCE(AVG(temperature), 0.0) AS temp_avg
		FROM readings
		WHERE cage_id IN (SELECT id FROM cages WHERE farm_id = $1) AND recorded_at >= NOW() - INTERVAL '7 days'
		GROUP BY cage_id`, farmID)
	doAvgMap := make(map[uuid.UUID]float64)
	tempAvgMap := make(map[uuid.UUID]float64)
	for _, a := range avgs {
		doAvgMap[a.CageID] = a.DOAvg
		tempAvgMap[a.CageID] = a.TempAvg
	}

	var items []models.CageComparisonItem
	for _, cage := range cages {
		cage.Location = models.GeoJSONPoint{
			Type:        "Point",
			Coordinates: []float64{cage.Longitude, cage.Latitude},
		}

		items = append(items, models.CageComparisonItem{
			Cage:          cage,
			LatestReading: readingMap[cage.ID],
			IncidentCount: incMap[cage.ID],
			AlertCount:    alertMap[cage.ID],
			RecentDOAvg:   doAvgMap[cage.ID],
			RecentTempAvg: tempAvgMap[cage.ID],
		})
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(items))
}

// GetFarmHistory handles GET /farms/:id/history?range=30d.
func (h *FarmHandler) GetFarmHistory(c *gin.Context) {
	farmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid farm ID parameter", 400))
		return
	}

	if !h.verifyFarmOwnership(c, farmID) {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Access denied: Farm does not belong to farmer", 403))
		return
	}

	var readings []models.Reading
	query := `
		SELECT r.id, r.cage_id, r.dissolved_oxygen, r.temperature, r.ph, r.turbidity, r.recorded_at, r.source
		FROM readings r
		JOIN cages c ON r.cage_id = c.id
		WHERE c.farm_id = $1 AND r.recorded_at >= NOW() - INTERVAL '30 days'
		ORDER BY r.recorded_at ASC`

	err = h.db.Select(&readings, query, farmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch farm history: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(readings))
}

// GetFarmIncidents handles GET /farms/:id/incidents.
func (h *FarmHandler) GetFarmIncidents(c *gin.Context) {
	farmID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid farm ID parameter", 400))
		return
	}

	if !h.verifyFarmOwnership(c, farmID) {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Access denied: Farm does not belong to farmer", 403))
		return
	}

	var incidents []models.Incident
	query := `
		SELECT i.id, i.cage_id, i.incident_type, i.description, i.detected_at, i.resolved_at, i.status
		FROM incidents i
		JOIN cages c ON i.cage_id = c.id
		WHERE c.farm_id = $1
		ORDER BY i.detected_at DESC`

	err = h.db.Select(&incidents, query, farmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch farm incidents: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(incidents))
}
