package models

import (
	"time"

	"github.com/google/uuid"
)

// Farm represents an aquaculture farm owned by a farmer.
type Farm struct {
	ID        uuid.UUID `db:"id" json:"id"`
	FarmerID  uuid.UUID `db:"farmer_id" json:"farmer_id"`
	Name      string    `db:"name" json:"name"`
	Region    string    `db:"region" json:"region"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// CreateFarmRequest payload for POST /farms.
type CreateFarmRequest struct {
	Name   string `json:"name" binding:"required"`
	Region string `json:"region" binding:"required"`
}

// FarmDashboardView provides an aggregate summary for Level 2 Farm Dashboard.
type FarmDashboardView struct {
	Farm          Farm                `json:"farm"`
	TotalCages    int                 `json:"total_cages"`
	ActiveCages   int                 `json:"active_cages"`
	OpenIncidents int                 `json:"open_incidents"`
	RecentAlerts  []Alert             `json:"recent_alerts"`
	CageStatuses  []CageStatusSummary `json:"cage_statuses"`
}

// CageStatusSummary summarizes a cage's current state on the dashboard.
type CageStatusSummary struct {
	CageID          uuid.UUID `json:"cage_id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	LatestDO        float64   `json:"latest_dissolved_oxygen"`
	LatestTemp      float64   `json:"latest_temperature"`
	LatestpH        float64   `json:"latest_ph"`
	LatestTurbidity float64   `json:"latest_turbidity"`
	LastRecordedAt  time.Time `json:"last_recorded_at"`
}
