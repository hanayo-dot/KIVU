package models

import (
	"time"

	"github.com/google/uuid"
)

// Alert represents a generated system warning (cage, farm, or region scope).
type Alert struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Scope        string    `db:"scope" json:"scope"`                 // 'cage' | 'farm' | 'region'
	RelatedID    string    `db:"related_id" json:"related_id"`       // cage_id, farm_id, or region_name
	Severity     string    `db:"severity" json:"severity"`           // 'info' | 'warning' | 'critical'
	Message      string    `db:"message" json:"message"`
	TriggeredAt  time.Time `db:"triggered_at" json:"triggered_at"`
	Acknowledged bool      `db:"acknowledged" json:"acknowledged"`
}
