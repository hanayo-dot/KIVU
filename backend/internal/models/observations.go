package models

import (
	"time"

	"github.com/google/uuid"
)

// Observation represents a qualitative observation logged by a farmer.
type Observation struct {
	ID              uuid.UUID `db:"id" json:"id"`
	CageID          uuid.UUID `db:"cage_id" json:"cage_id"`
	FarmerID        uuid.UUID `db:"farmer_id" json:"farmer_id"`
	ObservationType string    `db:"observation_type" json:"observation_type"` // 'fish_stress'|'mortality'|'discoloration'|'algal_growth'|'other'
	Description     string    `db:"description" json:"description"`
	Severity        string    `db:"severity" json:"severity"`                 // 'low'|'medium'|'high'
	RecordedAt      time.Time `db:"recorded_at" json:"recorded_at"`
}

// CreateObservationRequest payload.
type CreateObservationRequest struct {
	ObservationType string `json:"observation_type" binding:"required"`
	Description     string `json:"description" binding:"required"`
	Severity        string `json:"severity" binding:"required"`
}
