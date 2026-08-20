package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/hanayo-dot/KIVU/backend/config"
)

// DB wraps the sqlx.DB connection pool.
type DB struct {
	*sqlx.DB
}

// ConnectPostgres initializes a sqlx connection pool to PostgreSQL.
func ConnectPostgres(cfg *config.Config) (*DB, error) {
	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Printf("Warning: Database ping failed (%v). Ensure PostgreSQL is running with PostGIS.", err)
	} else {
		log.Println("Successfully connected to PostgreSQL (PostGIS ready)")
	}

	return &DB{DB: db}, nil
}
