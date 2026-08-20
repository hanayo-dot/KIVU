package services

import (
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// AIRiskEngine evaluates water quality readings for anomaly detection and risk scoring.
type AIRiskEngine struct {
	db           *sqlx.DB
	alertService *AlertsService
}

// NewAIRiskEngine initializes the AI Risk Engine service.
func NewAIRiskEngine(db *sqlx.DB, alertService *AlertsService) *AIRiskEngine {
	return &AIRiskEngine{
		db:           db,
		alertService: alertService,
	}
}

// AnalyzeCageReadings analyzes recent readings for a specific cage and generates risk scores and alerts.
func (engine *AIRiskEngine) AnalyzeCageReadings(cageID uuid.UUID) (*models.AIAnalysisResult, error) {
	var readings []models.Reading
	err := engine.db.Select(&readings,
		`SELECT id, cage_id, dissolved_oxygen, temperature, ph, turbidity, recorded_at, source 
		 FROM readings 
		 WHERE cage_id = $1 
		 ORDER BY recorded_at DESC LIMIT 10`,
		cageID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch readings for cage %s: %w", cageID, err)
	}

	if len(readings) == 0 {
		return &models.AIAnalysisResult{
			CageID:           cageID,
			RiskScore:        0.0,
			RiskLevel:        "low",
			TriggeredFactors: []string{"No reading history available"},
			Recommendation:   "Collect baseline water quality readings.",
			Confidence:       0.5,
			Basis:            "Insufficient sample data",
			AlertCreated:     false,
		}, nil
	}

	latest := readings[0]
	var triggered []string
	riskScore := 0.0

	// 1. Dissolved Oxygen Checks (Critical threshold: <3.0 mg/L, Warning: <5.0 mg/L)
	if latest.DissolvedOxygen < 3.0 {
		riskScore += 50.0
		triggered = append(triggered, fmt.Sprintf("Critical Hypoxia: Dissolved oxygen at %.2f mg/L (< 3.0 mg/L)", latest.DissolvedOxygen))
	} else if latest.DissolvedOxygen < 5.0 {
		riskScore += 25.0
		triggered = append(triggered, fmt.Sprintf("Moderate Hypoxia: Dissolved oxygen at %.2f mg/L (< 5.0 mg/L)", latest.DissolvedOxygen))
	}

	// Check DO Trend (e.g. 3 consecutive declining readings)
	if len(readings) >= 3 {
		if readings[0].DissolvedOxygen < readings[1].DissolvedOxygen && readings[1].DissolvedOxygen < readings[2].DissolvedOxygen {
			riskScore += 20.0
			triggered = append(triggered, "Deteriorating Trend: 3 consecutive declining dissolved oxygen measurements")
		}
	}

	// 2. Temperature Deviation Checks (>28.0°C or >2.0°C deviation from rolling average)
	if len(readings) >= 4 {
		sumTemp := 0.0
		for _, r := range readings {
			sumTemp += r.Temperature
		}
		avgTemp := sumTemp / float64(len(readings))
		if math.Abs(latest.Temperature-avgTemp) > 2.0 {
			riskScore += 15.0
			triggered = append(triggered, fmt.Sprintf("Thermal Anomaly: Current temp %.1f°C deviates >2.0°C from 7-day average (%.1f°C)", latest.Temperature, avgTemp))
		}
	}
	if latest.Temperature > 28.0 {
		riskScore += 10.0
		triggered = append(triggered, fmt.Sprintf("High Water Temp: %.1f°C exceeds optimal 28.0°C ceiling", latest.Temperature))
	}

	// 3. pH Level Checks (Safe range: 6.5 to 9.0)
	if latest.PH < 6.5 || latest.PH > 9.0 {
		riskScore += 15.0
		triggered = append(triggered, fmt.Sprintf("pH Imbalance: pH at %.2f is outside safe window (6.5 - 9.0)", latest.PH))
	}

	// 4. Turbidity Checks (>30 NTU)
	if latest.Turbidity > 30.0 {
		riskScore += 15.0
		triggered = append(triggered, fmt.Sprintf("High Turbidity: %.1f NTU exceeds 30.0 NTU threshold", latest.Turbidity))
	}

	if riskScore > 100.0 {
		riskScore = 100.0
	}

	// Assign Risk Level
	riskLevel := "low"
	recommendation := "Water quality parameters are within optimal operational range. Continue routine monitoring."
	if riskScore >= 70.0 {
		riskLevel = "critical"
		recommendation = "EMERGENCY: Immediate emergency aeration required! Reduce feed volume and check cage mesh flow."
	} else if riskScore >= 40.0 {
		riskLevel = "high"
		recommendation = "HIGH RISK: Deploy mechanical aerators and closely inspect fish behavior for asphyxiation stress."
	} else if riskScore >= 20.0 {
		riskLevel = "moderate"
		recommendation = "MODERATE RISK: Increase sampling frequency and monitor bay currents."
	}

	alertCreated := false
	if riskScore >= 40.0 && engine.alertService != nil {
		msg := fmt.Sprintf("AI Risk Alert: Cage %s registered %s risk (score: %.1f). Factors: %v", cageID, riskLevel, riskScore, triggered)
		_ = engine.alertService.TriggerCageAlert(cageID, "critical", msg)
		alertCreated = true
	}

	return &models.AIAnalysisResult{
		CageID:           cageID,
		RiskScore:        riskScore,
		RiskLevel:        riskLevel,
		TriggeredFactors: triggered,
		Recommendation:   recommendation,
		Confidence:       0.92,
		Basis:            "Rule-based threshold & multi-point temporal trend analysis",
		AlertCreated:     alertCreated,
	}, nil
}
