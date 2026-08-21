package services

import (
	"testing"

	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

func TestGroqServiceFallback(t *testing.T) {
	service := NewGroqService()

	telemetry := &models.CopernicusPointTelemetry{
		ZoneName:           "Winam Gulf (Kisumu Sector)",
		Latitude:           -0.45,
		Longitude:          34.20,
		DissolvedOxygen:    3.2, // Critical low DO
		SurfaceTemperature: 29.1,
		PH:                 8.2,
		Turbidity:          28.5,
		ChlorophyllA:       25.0,
		AlgalBloomRisk:     "high",
		RiskLevel:          "high",
		Trend:              "deteriorating",
	}

	rec, err := service.GenerateRecommendations(telemetry)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rec == nil {
		t.Fatalf("Expected recommendation object, got nil")
	}

	if len(rec.ActionableSteps) == 0 {
		t.Errorf("Expected actionable steps for critical DO telemetry, got 0")
	}

	if rec.RiskLevel != "critical" {
		t.Errorf("Expected critical risk level for DO=3.2 mg/L, got %s", rec.RiskLevel)
	}
}
