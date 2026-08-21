package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hanayo-dot/KIVU/backend/config"
)

// SMSClient interface for mocking and alternative providers
type SMSClient interface {
	SendSMS(to string, message string) error
}

type AfricasTalkingClient struct {
	cfg *config.Config
}

func NewSMSClient(cfg *config.Config) SMSClient {
	return &AfricasTalkingClient{cfg: cfg}
}

// SendSMS sends an SMS message via Africa's Talking API
func (c *AfricasTalkingClient) SendSMS(to string, message string) error {
	// For Hackathon MVP, if credentials aren't set, just log and pretend it sent.
	if c.cfg.Environment == "development" || c.cfg.Environment == "test" {
		log.Printf("[SMS] (Dev Mode) To: %s | Message: %s", to, message)
		return nil
	}

	// This is a placeholder for the actual Africa's Talking API call
	// URL: https://api.africastalking.com/version1/messaging
	apiUrl := "https://api.africastalking.com/version1/messaging"
	
	payload := map[string]string{
		"username": "sandbox", // use cfg.ATUsername
		"to":       to,
		"message":  message,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", "dummy_api_key") // use cfg.ATApiKey

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send SMS: HTTP %d", resp.StatusCode)
	}

	log.Printf("[SMS] Successfully sent alert to %s", to)
	return nil
}
