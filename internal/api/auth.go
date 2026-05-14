package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"aspera-terminal-ui/internal/config"
)

type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

func FetchToken(endpoint, username, password, clientID, clientSecret string) (*TokenResponse, error) {
	reqURL := fmt.Sprintf("%s/oauth/Customer/MOL/v1/token", strings.TrimRight(endpoint, "/"))

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)

	return doTokenRequest(reqURL, data, clientID, clientSecret)
}

func RefreshAccessToken(endpoint, refreshToken, clientID, clientSecret string) (*TokenResponse, error) {
	reqURL := fmt.Sprintf("%s/oauth/Customer/MOL/v1/token", strings.TrimRight(endpoint, "/"))

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	return doTokenRequest(reqURL, data, clientID, clientSecret)
}

func doTokenRequest(reqURL string, data url.Values, clientID, clientSecret string) (*TokenResponse, error) {
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch token, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var tResp TokenResponse
	if err := json.Unmarshal(bodyBytes, &tResp); err != nil {
		return nil, err
	}

	return &tResp, nil
}

// EnsureValidToken checks if the current token is valid and refreshes it if necessary.
func EnsureValidToken() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	if cfg.AccessToken == "" {
		return "", fmt.Errorf("not logged in")
	}

	// Check if token is expired (with 1 minute buffer)
	expiry, err := time.Parse(time.RFC3339, cfg.ExpiresAt)
	if err == nil && time.Now().Add(1*time.Minute).Before(expiry) {
		return cfg.AccessToken, nil
	}

	// Token is expired, try to refresh
	var tResp *TokenResponse
	
	// 1. Try refresh token if available
	if cfg.RefreshToken != "" {
		refreshExpiry, err := time.Parse(time.RFC3339, cfg.RefreshExpAt)
		if err == nil && time.Now().Before(refreshExpiry) {
			tResp, _ = RefreshAccessToken(cfg.Endpoint, cfg.RefreshToken, cfg.ClientID, cfg.ClientSecret)
		}
	}

	// 2. If refresh failed or not possible, try password re-login
	if tResp == nil && cfg.Username != "" {
		password, err := config.GetPassword(cfg.Username)
		if err == nil && password != "" {
			tResp, _ = FetchToken(cfg.Endpoint, cfg.Username, password, cfg.ClientID, cfg.ClientSecret)
		}
	}

	if tResp == nil {
		return "", fmt.Errorf("session expired, please login again")
	}

	// Update and save config
	cfg.AccessToken = tResp.AccessToken
	cfg.RefreshToken = tResp.RefreshToken
	cfg.ExpiresAt = time.Now().Add(time.Duration(tResp.ExpiresIn) * time.Second).Format(time.RFC3339)
	cfg.RefreshExpAt = time.Now().Add(time.Duration(tResp.RefreshTokenExpiresIn) * time.Second).Format(time.RFC3339)
	
	if err := config.SaveConfig(cfg); err != nil {
		return "", fmt.Errorf("failed to save refreshed token: %v", err)
	}

	return cfg.AccessToken, nil
}

