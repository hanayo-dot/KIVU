package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/hanayo-dot/KIVU/backend/internal/services"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// AIHandler manages AI Risk Engine endpoints.
type AIHandler struct {
	db       *sqlx.DB
	aiEngine *services.AIRiskEngine
}

// NewAIHandler initializes AIHandler.
func NewAIHandler(db *sqlx.DB, aiEngine *services.AIRiskEngine) *AIHandler {
	return &AIHandler{
		db:       db,
		aiEngine: aiEngine,
	}
}

// AnalyzeCage handles POST /ai/analyze-cage/:id.
func (h *AIHandler) AnalyzeCage(c *gin.Context) {
	cageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid cage ID parameter", 400))
		return
	}

	result, err := h.aiEngine.AnalyzeCageReadings(cageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("AI analysis failed: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(result))
}

// AnalyzeZone handles POST /ai/analyze-zone/:zone_id.
func (h *AIHandler) AnalyzeZone(c *gin.Context) {
	zoneID, err := uuid.Parse(c.Param("zone_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid zone ID parameter", 400))
		return
	}

	var metrics models.ZoneMetrics
	err = h.db.Get(&metrics, `SELECT id, zone_id, period, avg_dissolved_oxygen, avg_temperature, avg_ph, avg_turbidity, risk_level, trend, computed_at FROM zone_metrics WHERE zone_id = $1 ORDER BY computed_at DESC LIMIT 1`, zoneID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse("No zone metrics found to analyze", 404))
		return
	}

	riskLevel := "low"
	rationale := "Optimal environmental conditions across satellite observations."
	suitability := "high_suitability"

	if metrics.AvgDissolvedOxygen < 3.5 || metrics.AvgTurbidity > 40.0 {
		riskLevel = "high"
		suitability = "high_risk"
		rationale = fmt.Sprintf("Zone experiencing severe environmental stress: DO %.2f mg/L, Turbidity %.1f NTU", metrics.AvgDissolvedOxygen, metrics.AvgTurbidity)
	} else if metrics.AvgDissolvedOxygen < 5.0 || metrics.AvgTurbidity > 25.0 {
		riskLevel = "moderate"
		suitability = "watch"
		rationale = fmt.Sprintf("Zone on active watch: DO %.2f mg/L, Turbidity %.1f NTU", metrics.AvgDissolvedOxygen, metrics.AvgTurbidity)
	}

	// Update zone metrics and expansion signals
	_, _ = h.db.Exec(`UPDATE zone_metrics SET risk_level = $1 WHERE id = $2`, riskLevel, metrics.ID)
	_, _ = h.db.Exec(`INSERT INTO expansion_signals (zone_id, suitability, rationale) VALUES ($1, $2, $3)`, zoneID, suitability, rationale)

	type ZoneAnalysisResult struct {
		ZoneID      uuid.UUID `json:"zone_id"`
		RiskLevel   string    `json:"risk_level"`
		Suitability string    `json:"suitability"`
		Rationale   string    `json:"rationale"`
	}

	res := ZoneAnalysisResult{
		ZoneID:      zoneID,
		RiskLevel:   riskLevel,
		Suitability: suitability,
		Rationale:   rationale,
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(res))
}
