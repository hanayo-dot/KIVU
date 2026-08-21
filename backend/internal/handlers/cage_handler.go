package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/hanayo-dot/KIVU/backend/internal/middleware"
	"github.com/hanayo-dot/KIVU/backend/internal/services"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// CageHandler manages Level 1 My Cage endpoints.
type CageHandler struct {
	db           *sqlx.DB
	aiEngine     *services.AIRiskEngine
	alertService *services.AlertsService
}

// NewCageHandler initializes CageHandler.
func NewCageHandler(db *sqlx.DB, aiEngine *services.AIRiskEngine, alertService *services.AlertsService) *CageHandler {
	return &CageHandler{
		db:           db,
		aiEngine:     aiEngine,
		alertService: alertService,
	}
}

func (h *CageHandler) verifyCageOwnership(c *gin.Context, cageID uuid.UUID) bool {
	farmerID, ok := middleware.GetAuthenticatedFarmerID(c)
	if !ok {
		return false
	}
	var count int
	_ = h.db.Get(&count, `SELECT COUNT(*) FROM cages c JOIN farms f ON c.farm_id = f.id WHERE c.id = $1 AND f.farmer_id = $2`, cageID, farmerID)
	return count > 0
}

// SubmitReading handles POST /cages/:id/readings.
func (h *CageHandler) SubmitReading(c *gin.Context) {
	cageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid cage ID parameter", 400))
		return
	}

	if !h.verifyCageOwnership(c, cageID) {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Access denied: Cage does not belong to farmer", 403))
		return
	}

	var req models.CreateReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid reading payload: "+err.Error(), 400))
		return
	}

	// Telemetry boundary validation
	if req.DissolvedOxygen < 0.0 || req.DissolvedOxygen > 25.0 {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Dissolved oxygen out of valid range (0.0 - 25.0 mg/L)", 400))
		return
	}
	if req.Temperature < 0.0 || req.Temperature > 50.0 {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Temperature out of valid range (0.0 - 50.0 °C)", 400))
		return
	}
	if req.PH < 0.0 || req.PH > 14.0 {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("pH out of valid range (0.0 - 14.0)", 400))
		return
	}
	if req.Turbidity < 0.0 {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Turbidity must be a non-negative number", 400))
		return
	}

	source := req.Source
	if source == "" {
		source = "manual"
	}

	var reading models.Reading
	err = h.db.Get(&reading,
		`INSERT INTO readings (cage_id, dissolved_oxygen, temperature, ph, turbidity, source) 
		 VALUES ($1, $2, $3, $4, $5, $6) 
		 RETURNING id, cage_id, dissolved_oxygen, temperature, ph, turbidity, recorded_at, source`,
		cageID, req.DissolvedOxygen, req.Temperature, req.PH, req.Turbidity, source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to record reading: "+err.Error(), 500))
		return
	}

	// Auto-trigger AI Analysis & Alert escalation evaluation asynchronously
	go func() {
		if h.aiEngine != nil {
			_, _ = h.aiEngine.AnalyzeCageReadings(cageID)
		}
	}()

	c.JSON(http.StatusCreated, models.NewSuccessResponse(reading))
}

// GetLatestReading handles GET /cages/:id/readings/latest.
func (h *CageHandler) GetLatestReading(c *gin.Context) {
	cageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid cage ID parameter", 400))
		return
	}

	if !h.verifyCageOwnership(c, cageID) {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Access denied: Cage does not belong to farmer", 403))
		return
	}

	var reading models.Reading
	err = h.db.Get(&reading,
		`SELECT id, cage_id, dissolved_oxygen, temperature, ph, turbidity, recorded_at, source 
		 FROM readings 
		 WHERE cage_id = $1 
		 ORDER BY recorded_at DESC LIMIT 1`, cageID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse("No readings found for cage", 404))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(reading))
}

// GetReadingsHistory handles GET /cages/:id/readings?range=7d.
func (h *CageHandler) GetReadingsHistory(c *gin.Context) {
	cageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid cage ID parameter", 400))
		return
	}

	if !h.verifyCageOwnership(c, cageID) {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Access denied: Cage does not belong to farmer", 403))
		return
	}

	timeRange := c.DefaultQuery("range", "7d")
	interval := "7 days"
	if timeRange == "30d" {
		interval = "30 days"
	} else if timeRange == "24h" {
		interval = "24 hours"
	}

	var readings []models.Reading
	query := `SELECT id, cage_id, dissolved_oxygen, temperature, ph, turbidity, recorded_at, source 
			  FROM readings 
			  WHERE cage_id = $1 AND recorded_at >= NOW() - $2::interval 
			  ORDER BY recorded_at ASC`

	err = h.db.Select(&readings, query, cageID, interval)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch readings history: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(readings))
}

// SubmitObservation handles POST /cages/:id/observations.
func (h *CageHandler) SubmitObservation(c *gin.Context) {
	cageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid cage ID parameter", 400))
		return
	}

	farmerID, ok := middleware.GetAuthenticatedFarmerID(c)
	if !ok || !h.verifyCageOwnership(c, cageID) {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Access denied: Cage does not belong to farmer", 403))
		return
	}

	var req models.CreateObservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid observation payload: "+err.Error(), 400))
		return
	}

	var obs models.Observation
	err = h.db.Get(&obs,
		`INSERT INTO observations (cage_id, farmer_id, observation_type, description, severity) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id, cage_id, farmer_id, observation_type, description, severity, recorded_at`,
		cageID, farmerID, req.ObservationType, req.Description, req.Severity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to record observation: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(obs))
}

// GetCageAlerts handles GET /cages/:id/alerts.
func (h *CageHandler) GetCageAlerts(c *gin.Context) {
	cageIDStr := c.Param("id")
	cageID, err := uuid.Parse(cageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid cage ID parameter", 400))
		return
	}

	if !h.verifyCageOwnership(c, cageID) {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("Access denied: Cage does not belong to farmer", 403))
		return
	}

	var alerts []models.Alert
	err = h.db.Select(&alerts,
		`SELECT id, scope, related_id, severity, message, triggered_at, acknowledged 
		 FROM alerts 
		 WHERE scope = 'cage' AND related_id = $1 
		 ORDER BY triggered_at DESC`, cageIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch cage alerts: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(alerts))
}
