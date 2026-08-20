package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hanayo-dot/KIVU/backend/internal/services"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// AlertHandler manages system warnings and acknowledgments.
type AlertHandler struct {
	alertService *services.AlertsService
}

// NewAlertHandler initializes AlertHandler.
func NewAlertHandler(alertService *services.AlertsService) *AlertHandler {
	return &AlertHandler{alertService: alertService}
}

// GetAlerts handles GET /alerts?scope=cage|farm|region.
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	scope := c.Query("scope")
	alerts, err := h.alertService.GetAlerts(scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to fetch alerts: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(alerts))
}

// AcknowledgeAlert handles PATCH /alerts/:id/acknowledge.
func (h *AlertHandler) AcknowledgeAlert(c *gin.Context) {
	alertID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid alert ID parameter", 400))
		return
	}

	err = h.alertService.AcknowledgeAlert(alertID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(err.Error(), 404))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
		"message": "Alert acknowledged successfully",
		"id":      alertID,
	}))
}
