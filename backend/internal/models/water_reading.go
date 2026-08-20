package models

import (
	"time"

	"github.com/google/uuid"
)

// Reading represents a water-quality measurement entry for a cage.
type Reading struct {
	ID              uuid.UUID `db:"id" json:"id"`
	CageID          uuid.UUID `db:"cage_id" json:"cage_id"`
	DissolvedOxygen float64   `db:"dissolved_oxygen" json:"dissolved_oxygen"` // mg/L
	Temperature     float64   `db:"temperature" json:"temperature"`           // °C
	PH              float64   `db:"ph" json:"ph"`
	Turbidity       float64   `db:"turbidity" json:"turbidity"`               // NTU
	RecordedAt      time.Time `db:"recorded_at" json:"recorded_at"`
	Source          string    `db:"source" json:"source"`                     // 'sensor' | 'manual'
}

// CreateReadingRequest payload.
type CreateReadingRequest struct {
	DissolvedOxygen float64 `json:"dissolved_oxygen" binding:"required"`
	Temperature     float64 `json:"temperature" binding:"required"`
	PH              float64 `json:"ph" binding:"required"`
	Turbidity       float64 `json:"turbidity" binding:"required"`
	Source          string  `json:"source"` // default 'manual' or 'sensor'
}
