package integrations

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hanayo-dot/KIVU/backend/config"
)

func TestCopernicusOAuthTokenFetching(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}
		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("expected grant_type=client_credentials, got %s", r.FormValue("grant_type"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "mock_copernicus_access_token_12345",
			"expires_in":   3600,
			"token_type":   "Bearer",
		})
	}))
	defer server.Close()

	cfg := &config.Config{
		CopernicusClientID:     "test_id",
		CopernicusClientSecret: "test_secret",
	}

	client := NewCopernicusClient(cfg)
	client.TokenURL = server.URL

	token, err := client.GetAccessToken()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if token != "mock_copernicus_access_token_12345" {
		t.Errorf("expected token mock_copernicus_access_token_12345, got %s", token)
	}
}

func TestCopernicusTokenCaching(t *testing.T) {
	cfg := &config.Config{
		CopernicusClientID:     "test_id",
		CopernicusClientSecret: "test_secret",
	}

	client := NewCopernicusClient(cfg)
	client.accessToken = "cached_test_token"
	client.expiresAt = time.Now().Add(30 * time.Minute)

	token, err := client.GetAccessToken()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token != "cached_test_token" {
		t.Errorf("expected cached token to be reused, got %s", token)
	}
}
