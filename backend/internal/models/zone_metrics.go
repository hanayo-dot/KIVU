package models

import (
	"time"

	"github.com/google/uuid"
)

// ZoneMetrics stores aggregated or Copernicus satellite-derived environmental data per zone.
type ZoneMetrics struct {
	ID                 uuid.UUID `db:"id" json:"id"`
	ZoneID             uuid.UUID `db:"zone_id" json:"zone_id"`
	Period             string    `db:"period" json:"period"`
	AvgDissolvedOxygen float64   `db:"avg_dissolved_oxygen" json:"avg_dissolved_oxygen"`
	AvgTemperature     float64   `db:"avg_temperature" json:"avg_temperature"`
	AvgPH              float64   `db:"avg_ph" json:"avg_ph"`
	AvgTurbidity       float64   `db:"avg_turbidity" json:"avg_turbidity"`
	RiskLevel          string    `db:"risk_level" json:"risk_level"` // 'low' | 'moderate' | 'high'
	Trend              string    `db:"trend" json:"trend"`           // 'improving' | 'stable' | 'deteriorating'
	ComputedAt         time.Time `db:"computed_at" json:"computed_at"`
}
