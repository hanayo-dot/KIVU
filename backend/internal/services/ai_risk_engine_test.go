package services

import (
	"testing"
)

// RiskHeuristics evaluates water telemetry parameter thresholds.
func evaluateParameterRisk(dissolvedOxygen, temp, ph, ammonia float64) (float64, string) {
	riskScore := 0.0
	details := ""

	// Dissolved Oxygen critical (< 5.0 mg/L)
	if dissolvedOxygen < 3.0 {
		riskScore += 50.0
		details += "Critical hypoxia (DO < 3.0 mg/L). "
	} else if dissolvedOxygen < 5.0 {
		riskScore += 25.0
		details += "Low dissolved oxygen (DO < 5.0 mg/L). "
	}

	// Temperature anomalies (> 29°C or < 22°C)
	if temp > 30.0 || temp < 20.0 {
		riskScore += 30.0
		details += "Severe water temperature stress. "
	} else if temp > 28.0 || temp < 22.0 {
		riskScore += 15.0
		details += "Moderate temperature anomaly. "
	}

	// pH levels (< 6.5 or > 8.5)
	if ph < 6.5 || ph > 8.5 {
		riskScore += 15.0
		details += "Abnormal pH level. "
	}

	// Ammonia concentration (> 0.05 mg/L)
	if ammonia > 0.05 {
		riskScore += 20.0
		details += "High un-ionized ammonia concentration. "
	}

	if riskScore > 100.0 {
		riskScore = 100.0
	}
	return riskScore, details
}

func TestAIRiskEngineParameterEvaluation(t *testing.T) {
	tests := []struct {
		name          string
		do            float64
		temp          float64
		ph            float64
		ammonia       float64
		minExpected   float64
		expectAnomaly bool
	}{
		{
			name:          "Optimal Water Conditions",
			do:            7.5,
			temp:          25.0,
			ph:            7.2,
			ammonia:       0.01,
			minExpected:   0.0,
			expectAnomaly: false,
		},
		{
			name:          "Hypoxia Low Oxygen Event",
			do:            2.5,
			temp:          25.0,
			ph:            7.2,
			ammonia:       0.01,
			minExpected:   50.0,
			expectAnomaly: true,
		},
		{
			name:          "High Temp and High Ammonia Compound Risk",
			do:            4.5,
			temp:          31.0,
			ph:            8.8,
			ammonia:       0.08,
			minExpected:   80.0,
			expectAnomaly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := evaluateParameterRisk(tt.do, tt.temp, tt.ph, tt.ammonia)
			if score < tt.minExpected {
				t.Errorf("%s: expected min risk score %v, got %v", tt.name, tt.minExpected, score)
			}
			isAnomaly := score >= 25.0
			if isAnomaly != tt.expectAnomaly {
				t.Errorf("%s: expected anomaly state %v, got %v", tt.name, tt.expectAnomaly, isAnomaly)
			}
		})
	}
}
