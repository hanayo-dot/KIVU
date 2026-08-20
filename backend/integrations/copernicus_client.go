package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/hanayo-dot/KIVU/backend/config"
)

// CopernicusClient manages CDSE OAuth2 token fetching and Statistical API requests.
type CopernicusClient struct {
	cfg         *config.Config
	accessToken string
	expiresAt   time.Time
	mu          sync.Mutex
}

// NewCopernicusClient initializes a new CopernicusClient.
func NewCopernicusClient(cfg *config.Config) *CopernicusClient {
	return &CopernicusClient{cfg: cfg}
}

// GetAccessToken fetches or returns a cached valid OAuth2 token using client credentials grant.
func (c *CopernicusClient) GetAccessToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.expiresAt.Add(-1*time.Minute)) {
		return c.accessToken, nil
	}

	if c.cfg.CopernicusClientID == "" || c.cfg.CopernicusClientSecret == "" {
		return "", fmt.Errorf("copernicus credentials missing in environment")
	}

	tokenURL := "https://identity.dataspace.copernicus.eu/auth/realms/CDSE/protocol/openid-connect/token"
	resp, err := http.PostForm(tokenURL, url.Values{
		"client_id":     {c.cfg.CopernicusClientID},
		"client_secret": {c.cfg.CopernicusClientSecret},
		"grant_type":    {"client_credentials"},
	})

	if err != nil {
		return "", fmt.Errorf("oauth token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oauth request returned non-200 status: %d", resp.StatusCode)
	}

	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("failed to decode oauth response: %w", err)
	}

	c.accessToken = body.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(body.ExpiresIn) * time.Second)

	return c.accessToken, nil
}
