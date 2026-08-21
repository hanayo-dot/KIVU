package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/hanayo-dot/KIVU/backend/config"
	"github.com/hanayo-dot/KIVU/backend/integrations"
	"github.com/hanayo-dot/KIVU/backend/internal/handlers"
	"github.com/hanayo-dot/KIVU/backend/internal/middleware"
	"github.com/hanayo-dot/KIVU/backend/internal/services"
)

// SetupRouter initializes Gin engine with CORS, logging, JWT Auth middleware, and all KIVU REST endpoints.
func SetupRouter(cfg *config.Config, db *sqlx.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLogger())

	// Initialize Services
	authService := services.NewAuthService(cfg, db)
	smsClient := integrations.NewSMSClient(cfg)
	alertService := services.NewAlertsService(db, smsClient)
	aiEngine := services.NewAIRiskEngine(db, alertService)
	expansionService := services.NewExpansionService(db)

	// Start Copernicus background sync
	copClient := integrations.NewCopernicusClient(cfg)
	copSync := services.NewCopernicusSyncService(cfg, db, copClient)
	copSync.StartSyncLoop()

	// Initialize Handlers
	authH := handlers.NewAuthHandler(authService)
	onboardingH := handlers.NewOnboardingHandler(db)
	cageH := handlers.NewCageHandler(db, aiEngine, alertService)
	farmH := handlers.NewFarmHandler(db)
	lakeH := handlers.NewLakeHandler(db)
	expansionH := handlers.NewExpansionHandler(expansionService)
	aiH := handlers.NewAIHandler(db, aiEngine)
	alertH := handlers.NewAlertHandler(alertService)

	// Public Endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": "KIVU Backend MVP with Auth"})
	})

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authH.Register)
		authGroup.POST("/login", authH.Login)
		authGroup.POST("/refresh", authH.Refresh)
		authGroup.POST("/logout", authH.Logout)
	}

	// Public Spatial Lake Map View
	r.GET("/lake/zones", lakeH.GetLakeZones)

	// Protected Endpoints (Require Bearer JWT)
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired(cfg))
	{
		// Onboarding
		protected.POST("/farms", onboardingH.RegisterFarm)
		protected.POST("/cages", onboardingH.RegisterCage)

		// Level 1 — My Cage
		protected.POST("/cages/:id/readings", cageH.SubmitReading)
		protected.GET("/cages/:id/readings/latest", cageH.GetLatestReading)
		protected.GET("/cages/:id/readings", cageH.GetReadingsHistory)
		protected.POST("/cages/:id/observations", cageH.SubmitObservation)
		protected.GET("/cages/:id/alerts", cageH.GetCageAlerts)

		// Level 2 — My Farm
		protected.GET("/farms/:id/dashboard", farmH.GetFarmDashboard)
		protected.GET("/farms/:id/cages/compare", farmH.CompareCages)
		protected.GET("/farms/:id/history", farmH.GetFarmHistory)
		protected.GET("/farms/:id/incidents", farmH.GetFarmIncidents)

		// Level 3 — Lake Victoria Detail & Alerts
		protected.GET("/lake/zones/:id", lakeH.GetZoneByID)
		protected.GET("/lake/alerts", lakeH.GetRegionalAlerts)
		protected.GET("/lake/hotspots", lakeH.GetHotspots)

		// Level 4 — Future Expansion Intelligence
		protected.GET("/expansion/signals", expansionH.GetAllSignals)
		protected.GET("/expansion/signals/:zone_id", expansionH.GetSignalByZone)
		protected.POST("/expansion/evaluate", expansionH.EvaluateCoordinates)

		// AI Risk Engine
		protected.POST("/ai/analyze-cage/:id", aiH.AnalyzeCage)
		protected.POST("/ai/analyze-zone/:zone_id", aiH.AnalyzeZone)

		// Alerts Management
		protected.GET("/alerts", alertH.GetAlerts)
		protected.PATCH("/alerts/:id/acknowledge", alertH.AcknowledgeAlert)
	}

	return r
}
