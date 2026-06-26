package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"aspera-terminal-ui/internal/config"
)

func TestFetchToken(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Verify request
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		if req.URL.Path != "/oauth/Customer/MOL/v1/token" {
			t.Errorf("expected /oauth/Customer/MOL/v1/token, got %s", req.URL.Path)
		}

		username, password, ok := req.BasicAuth()
		if !ok || username != "client-id" || password != "client-secret" {
			t.Errorf("invalid basic auth")
		}

		err := req.ParseForm()
		if err != nil {
			t.Fatal(err)
		}
		if req.PostFormValue("grant_type") != "password" ||
			req.PostFormValue("username") != "testuser" ||
			req.PostFormValue("password") != "testpass" {
			t.Errorf("invalid post form values: %v", req.PostForm)
		}

		resp := TokenResponse{
			AccessToken:           "new-access-token",
			RefreshToken:          "new-refresh-token",
			ExpiresIn:             3600,
			RefreshTokenExpiresIn: 7200,
		}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(resp)
	}))
	defer server.Close()

	tResp, err := FetchToken(server.URL, "testuser", "testpass", "client-id", "client-secret")
	if err != nil {
		t.Fatalf("FetchToken failed: %v", err)
	}

	if tResp.AccessToken != "new-access-token" || tResp.RefreshToken != "new-refresh-token" {
		t.Errorf("unexpected token response: %+v", tResp)
	}
}

func TestRefreshAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			t.Fatal(err)
		}
		if req.PostFormValue("grant_type") != "refresh_token" ||
			req.PostFormValue("refresh_token") != "old-refresh-token" {
			t.Errorf("invalid post form values: %v", req.PostForm)
		}

		resp := TokenResponse{
			AccessToken:           "refreshed-access-token",
			RefreshToken:          "new-refresh-token",
			ExpiresIn:             3600,
			RefreshTokenExpiresIn: 7200,
		}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(resp)
	}))
	defer server.Close()

	tResp, err := RefreshAccessToken(server.URL, "old-refresh-token", "client-id", "client-secret")
	if err != nil {
		t.Fatalf("RefreshAccessToken failed: %v", err)
	}

	if tResp.AccessToken != "refreshed-access-token" {
		t.Errorf("unexpected token response: %+v", tResp)
	}
}

func TestEnsureValidToken_Valid(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "api-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Save a valid config
	expiry := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	cfg := &config.Config{
		AccessToken: "valid-token",
		ExpiresAt:   expiry,
	}
	if err := config.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}

	tok, err := EnsureValidToken()
	if err != nil {
		t.Fatalf("EnsureValidToken failed: %v", err)
	}

	if tok != "valid-token" {
		t.Errorf("expected 'valid-token', got %q", tok)
	}
}

func TestEnsureValidToken_ExpiredRefresh(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "api-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		resp := TokenResponse{
			AccessToken:           "refreshed-token",
			RefreshToken:          "new-refresh-token",
			ExpiresIn:             3600,
			RefreshTokenExpiresIn: 7200,
		}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(resp)
	}))
	defer server.Close()

	// Save an expired config with valid refresh token
	cfg := &config.Config{
		AccessToken:  "expired-token",
		ExpiresAt:    time.Now().Add(-10 * time.Minute).Format(time.RFC3339),
		RefreshToken: "refresh-token",
		RefreshExpAt: time.Now().Add(10 * time.Minute).Format(time.RFC3339),
		Endpoint:     server.URL,
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}
	if err := config.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}

	tok, err := EnsureValidToken()
	if err != nil {
		t.Fatalf("EnsureValidToken failed: %v", err)
	}

	if tok != "refreshed-token" {
		t.Errorf("expected 'refreshed-token', got %q", tok)
	}

	// Verify it was saved to the config
	saved, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if saved.AccessToken != "refreshed-token" {
		t.Errorf("expected saved token to be 'refreshed-token', got %q", saved.AccessToken)
	}
}
