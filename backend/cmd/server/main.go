package main

import (
	"fmt"
	"log"

	"github.com/hanayo-dot/KIVU/backend/config"
	"github.com/hanayo-dot/KIVU/backend/database"
	"github.com/hanayo-dot/KIVU/backend/routes"
)

func main() {
	log.Println("==================================================")
	log.Println("  KIVU Aquaculture & Lake Intelligence Backend    ")
	log.Println("==================================================")

	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Fatal configuration error: %v", err)
	}
	log.Printf("Starting KIVU server on port :%s [env: %s]", cfg.ServerPort, cfg.Environment)

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("Fatal: Database initialization error: %v", err)
	}
	defer db.Close()

	router := routes.SetupRouter(cfg, db.DB)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("KIVU Backend ready and listening at http://localhost%s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Fatal: HTTP Server crashed: %v", err)
	}
}
