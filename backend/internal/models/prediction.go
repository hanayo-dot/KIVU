package models

import (
	"time"

	"github.com/google/uuid"
)

// ExpansionSignal classifies zone suitability for future cage expansion.
type ExpansionSignal struct {
	ID          uuid.UUID `db:"id" json:"id"`
	ZoneID      uuid.UUID `db:"zone_id" json:"zone_id"`
	Suitability string    `db:"suitability" json:"suitability"` // 'high_suitability' | 'watch' | 'high_risk'
	Rationale   string    `db:"rationale" json:"rationale"`
	ComputedAt  time.Time `db:"computed_at" json:"computed_at"`
}

// EvaluateCoordinatesRequest payload for POST /expansion/evaluate.
type EvaluateCoordinatesRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// AIAnalysisResult response structure for POST /ai/analyze-cage/:id.
type AIAnalysisResult struct {
	CageID           uuid.UUID `json:"cage_id"`
	RiskScore        float64   `json:"risk_score"`        // 0.0 to 100.0
	RiskLevel        string    `json:"risk_level"`        // 'low' | 'moderate' | 'high' | 'critical'
	TriggeredFactors []string  `json:"triggered_factors"`
	Recommendation   string    `json:"recommendation"`
	Confidence       float64   `json:"confidence"`
	Basis            string    `json:"basis"`
	AlertCreated     bool      `json:"alert_created"`
}
