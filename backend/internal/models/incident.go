package models

import (
	"time"

	"github.com/google/uuid"
)

// Incident represents a logged or AI-detected event.
type Incident struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	CageID       uuid.UUID  `db:"cage_id" json:"cage_id"`
	IncidentType string     `db:"incident_type" json:"incident_type"`
	Description  string     `db:"description" json:"description"`
	DetectedAt   time.Time  `db:"detected_at" json:"detected_at"`
	ResolvedAt   *time.Time `db:"resolved_at" json:"resolved_at,omitempty"`
	Status       string     `db:"status" json:"status"` // 'open' | 'resolved'
}
