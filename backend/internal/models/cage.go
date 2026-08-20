package models

import (
	"time"

	"github.com/google/uuid"
)

// GeoJSONPoint represents a GeoJSON Point object for frontend map integration.
type GeoJSONPoint struct {
	Type        string    `json:"type"`        // Always "Point"
	Coordinates []float64 `json:"coordinates"` // [longitude, latitude]
}

// Cage represents an individual fish cage on the lake.
type Cage struct {
	ID          uuid.UUID    `db:"id" json:"id"`
	FarmID      uuid.UUID    `db:"farm_id" json:"farm_id"`
	Name        string       `db:"name" json:"name"`
	Latitude    float64      `db:"latitude" json:"latitude"`
	Longitude   float64      `db:"longitude" json:"longitude"`
	Location    GeoJSONPoint `json:"location"`
	InstalledAt time.Time    `db:"installed_at" json:"installed_at"`
	Status      string       `db:"status" json:"status"`
}

// CreateCageRequest payload for POST /cages.
type CreateCageRequest struct {
	FarmID    uuid.UUID `json:"farm_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Latitude  float64   `json:"latitude" binding:"required"`
	Longitude float64   `json:"longitude" binding:"required"`
}

// CageComparisonItem for side-by-side comparison.
type CageComparisonItem struct {
	Cage          Cage      `json:"cage"`
	LatestReading *Reading  `json:"latest_reading"`
	IncidentCount int       `json:"incident_count"`
	AlertCount    int       `json:"alert_count"`
	RecentDOAvg   float64   `json:"recent_do_avg"`
	RecentTempAvg float64   `json:"recent_temp_avg"`
}
