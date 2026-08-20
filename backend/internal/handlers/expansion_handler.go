package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hanayo-dot/KIVU/backend/internal/services"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// ExpansionHandler manages Level 4 Future Expansion suitability endpoints.
type ExpansionHandler struct {
	expansionService *services.ExpansionService
}

// NewExpansionHandler initializes ExpansionHandler.
func NewExpansionHandler(expansionService *services.ExpansionService) *ExpansionHandler {
	return &ExpansionHandler{expansionService: expansionService}
}

// GetAllSignals handles GET /expansion/signals.
func (h *ExpansionHandler) GetAllSignals(c *gin.Context) {
	signals, err := h.expansionService.GetAllSignals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch expansion signals: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(signals))
}

// GetSignalByZone handles GET /expansion/signals/:zone_id.
func (h *ExpansionHandler) GetSignalByZone(c *gin.Context) {
	zoneID, err := uuid.Parse(c.Param("zone_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid zone ID parameter", 400))
		return
	}

	signal, err := h.expansionService.GetSignalByZone(zoneID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(err.Error(), 404))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(signal))
}

// EvaluateCoordinates handles POST /expansion/evaluate.
func (h *ExpansionHandler) EvaluateCoordinates(c *gin.Context) {
	var req models.EvaluateCoordinatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid coordinate payload: "+err.Error(), 400))
		return
	}

	signal, err := h.expansionService.EvaluateCoordinates(req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Spatial evaluation failed: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(signal))
}
