package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/hanayo-dot/KIVU/backend/internal/middleware"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// OnboardingHandler manages farm and cage creation for authenticated farmers.
type OnboardingHandler struct {
	db *sqlx.DB
}

// NewOnboardingHandler initializes OnboardingHandler.
func NewOnboardingHandler(db *sqlx.DB) *OnboardingHandler {
	return &OnboardingHandler{db: db}
}

// RegisterFarm handles POST /farms. Uses authenticated farmer_id from context.
func (h *OnboardingHandler) RegisterFarm(c *gin.Context) {
	farmerID, ok := middleware.GetAuthenticatedFarmerID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Unauthenticated", 401))
		return
	}

	var req models.CreateFarmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request payload: "+err.Error(), 400))
		return
	}

	var farm models.Farm
	err := h.db.Get(&farm,
		`INSERT INTO farms (farmer_id, name, region) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, farmer_id, name, region, created_at`,
		farmerID, req.Name, req.Region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to register farm: "+err.Error(), 500))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(farm))
}

// RegisterCage handles POST /cages. Verifies farm ownership by authenticated farmer.
func (h *OnboardingHandler) RegisterCage(c *gin.Context) {
	farmerID, ok := middleware.GetAuthenticatedFarmerID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Unauthenticated", 401))
		return
	}

	var req models.CreateCageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request payload: "+err.Error(), 400))
		return
	}

	// Verify farm ownership
	var count int
	_ = h.db.Get(&count, `SELECT COUNT(*) FROM farms WHERE id = $1 AND farmer_id = $2`, req.FarmID, farmerID)
	if count == 0 {
		c.JSON(http.StatusForbidden, models.NewErrorResponse("You do not own the target farm for this cage", 403))
		return
	}

	var cageID string
	query := `
		INSERT INTO cages (farm_id, name, latitude, longitude, location, status) 
		VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($4, $3), 4326)::geography, 'active') 
		RETURNING id::text`

	err := h.db.Get(&cageID, query, req.FarmID, req.Name, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to register cage: "+err.Error(), 500))
		return
	}

	res := models.Cage{
		ID:        models.MustParseUUID(cageID),
		FarmID:    req.FarmID,
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Location: models.GeoJSONPoint{
			Type:        "Point",
			Coordinates: []float64{req.Longitude, req.Latitude},
		},
		Status: "active",
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(res))
}
