package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hanayo-dot/KIVU/backend/internal/models"
)

// GroqService provides AI recommendations and actionable steps using the Groq API.
type GroqService struct {
	apiKey     string
	modelName  string
	httpClient *http.Client
}

// NewGroqService initializes a new Groq API client service.
func NewGroqService() *GroqService {
	apiKey := os.Getenv("GROQ_API_KEY")
	modelName := os.Getenv("GROQ_MODEL")
	if modelName == "" {
		modelName = "llama-3.3-70b-versatile"
	}

	return &GroqService{
		apiKey:    apiKey,
		modelName: modelName,
		httpClient: &http.Client{
			Timeout: 12 * time.Second,
		},
	}
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model          string        `json:"model"`
	Messages       []groqMessage `json:"messages"`
	ResponseFormat *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
	Temperature float64 `json:"temperature"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// GenerateRecommendations produces AI-driven actionable steps and recommendations for given telemetry.
func (s *GroqService) GenerateRecommendations(telemetry *models.CopernicusPointTelemetry) (*models.AIRecommendation, error) {
	if s.apiKey == "" {
		// Fallback to local expert rule synthesis when GROQ_API_KEY is not set
		return s.GenerateFallbackRecommendations(telemetry), nil
	}

	systemPrompt := `You are KIVU AI — an expert Lake Victoria Aquaculture Specialist and Copernicus Satellite Limnology Advisor.
Analyze the provided Copernicus satellite telemetry and output a strict JSON object with the following schema:
{
  "summary": "Short 1-2 sentence overall summary of water health and fish survival impact",
  "risk_score": 0.0 to 100.0 (number),
  "risk_level": "low" | "moderate" | "high" | "critical",
  "actionable_steps": ["Step 1...", "Step 2...", "Step 3..."],
  "preventative_measures": ["Measure 1...", "Measure 2..."]
}`

	userPrompt := fmt.Sprintf(`Copernicus Telemetry Data:
Location: %s (Lat: %.4f, Lng: %.4f)
Dissolved Oxygen: %.2f mg/L
Surface Temperature: %.1f °C
pH: %.2f
Turbidity: %.1f NTU
Chlorophyll-a: %.1f µg/L
Secchi Transparency: %.2f m
Algal Bloom Risk: %s
Current System Risk Level: %s
Trend: %s`,
		telemetry.ZoneName, telemetry.Latitude, telemetry.Longitude,
		telemetry.DissolvedOxygen, telemetry.SurfaceTemperature, telemetry.PH,
		telemetry.Turbidity, telemetry.ChlorophyllA, telemetry.SecchiDepth,
		telemetry.AlgalBloomRisk, telemetry.RiskLevel, telemetry.Trend)

	reqPayload := groqRequest{
		Model: s.modelName,
		Messages: []groqMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		ResponseFormat: &struct {
			Type string `json:"type"`
		}{Type: "json_object"},
		Temperature: 0.2,
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return s.GenerateFallbackRecommendations(telemetry), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return s.GenerateFallbackRecommendations(telemetry), nil
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		return s.GenerateFallbackRecommendations(telemetry), nil
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.GenerateFallbackRecommendations(telemetry), nil
	}

	var groqResp groqResponse
	if err := json.Unmarshal(respBytes, &groqResp); err != nil || len(groqResp.Choices) == 0 {
		return s.GenerateFallbackRecommendations(telemetry), nil
	}

	content := groqResp.Choices[0].Message.Content
	var rec models.AIRecommendation
	if err := json.Unmarshal([]byte(content), &rec); err != nil {
		return s.GenerateFallbackRecommendations(telemetry), nil
	}

	rec.AIModelUsed = fmt.Sprintf("%s (Groq AI API)", s.modelName)
	return &rec, nil
}

// GenerateFallbackRecommendations synthesizes expert actionable steps locally if Groq API is offline or unconfigured.
func (s *GroqService) GenerateFallbackRecommendations(telemetry *models.CopernicusPointTelemetry) *models.AIRecommendation {
	var steps []string
	var preventative []string
	riskScore := 20.0
	riskLevel := "low"
	summary := "Water conditions in this sector are optimal for aquaculture operations."

	// 1. Dissolved Oxygen Interventions
	if telemetry.DissolvedOxygen < 3.5 {
		riskScore += 50
		riskLevel = "critical"
		summary = "CRITICAL HYPOXIA DETECTED: Dissolved oxygen levels are severely low, posing immediate fish mortality risk."
		steps = append(steps, "CRITICAL: Deploy emergency mechanical aerators or oxygen diffusers immediately.")
		steps = append(steps, "Pause all feeding schedules to prevent oxygen depletion during fish digestion.")
		steps = append(steps, "Inspect net mesh for clogging to maximize water exchange through cages.")
	} else if telemetry.DissolvedOxygen < 5.0 {
		riskScore += 25
		riskLevel = "moderate"
		summary = "MODERATE HYPOXIA WARNING: Dissolved oxygen is below optimal levels."
		steps = append(steps, "Prepare backup paddlewheel aerators for evening deployment.")
		steps = append(steps, "Reduce feeding rations by 50% until oxygen levels recover > 5.5 mg/L.")
	} else {
		steps = append(steps, "Maintain regular feeding schedule during peak morning oxygen hours.")
	}

	// 2. Chlorophyll-a / Algal Bloom Interventions
	if telemetry.ChlorophyllA > 20.0 || telemetry.AlgalBloomRisk == "high" {
		riskScore += 20
		if riskLevel == "low" {
			riskLevel = "moderate"
		}
		steps = append(steps, "Monitor for microcystin cyanobacteria scum formation on the water surface.")
		steps = append(steps, "Avoid harvesting during active bloom peaks to preserve fish flesh quality.")
		preventative = append(preventative, "Install perimeter bubble curtains to deflect floating algal blooms from entering cages.")
	}

	// 3. Turbidity Interventions
	if telemetry.Turbidity > 25.0 {
		riskScore += 15
		steps = append(steps, "Check fish gills for silt clogging and signs of respiratory stress.")
		preventative = append(preventative, "Consider relocating cages further offshore into deeper, lower-turbidity waters.")
	}

	// 4. Surface Temperature Interventions
	if telemetry.SurfaceTemperature > 28.5 {
		riskScore += 15
		steps = append(steps, "Sink cages 0.5 - 1.0 meter deeper if adjustable to reach cooler water strata.")
	}

	if len(preventative) == 0 {
		preventative = append(preventative, "Maintain continuous Copernicus satellite monitoring for early warning signals.")
		preventative = append(preventative, "Perform weekly water quality calibration check with handheld probes.")
	}

	if riskScore > 100 {
		riskScore = 100
	}

	return &models.AIRecommendation{
		Summary:              summary,
		RiskScore:            riskScore,
		RiskLevel:            riskLevel,
		ActionableSteps:      steps,
		PreventativeMeasures: preventative,
		AIModelUsed:          "KIVU Local Rule-Engine (Groq Fallback)",
	}
}
