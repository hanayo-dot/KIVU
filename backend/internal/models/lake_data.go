package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// GeoJSONPolygon represents a GeoJSON Polygon object.
type GeoJSONPolygon struct {
	Type        string        `json:"type"` // Always "Polygon"
	Coordinates [][][]float64 `json:"coordinates"`
}

// LakeZone represents a spatial region of Lake Victoria.
type LakeZone struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	Name        string          `db:"name" json:"name"`
	Boundary    json.RawMessage `json:"boundary"` // Stored/Returned as GeoJSON Polygon
	RegionLabel string          `db:"region_label" json:"region_label"`
}

// LakeZoneDetail includes historical metrics, incident density, and active cages.
type LakeZoneDetail struct {
	Zone            LakeZone        `json:"zone"`
	CurrentMetrics  *ZoneMetrics    `json:"current_metrics"`
	ExpansionSignal *ExpansionSignal`json:"expansion_signal,omitempty"`
	IncidentCount   int             `json:"incident_count"`
	ActiveCageCount int             `json:"active_cage_count"`
}
