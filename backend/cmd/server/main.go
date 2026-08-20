package main

import (
	"log"

	"github.com/hanayo-dot/KIVU/backend/config"
	"github.com/hanayo-dot/KIVU/backend/database"
)

func main() {
	log.Println("==================================================")
	log.Println("  KIVU Aquaculture & Lake Intelligence Backend    ")
	log.Println("==================================================")

	cfg := config.LoadConfig()
	log.Printf("Starting KIVU server on port :%s [env: %s]", cfg.ServerPort, cfg.Environment)

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("Fatal: Database initialization error: %v", err)
	}
	defer db.Close()

	log.Println("KIVU Backend Core Foundation & PostGIS DB pool ready.")
}
