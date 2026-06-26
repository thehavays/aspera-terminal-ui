package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	expected := filepath.Join(tmpDir, "aspera-terminal-ui", "config.json")
	actual, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := &Config{
		AccessToken:  "access-123",
		RefreshToken: "refresh-123",
		ExpiresAt:    "expiry-time",
		RefreshExpAt: "refresh-expiry-time",
		Endpoint:     "https://test.endpoint",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Username:     "username",
	}

	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.AccessToken != cfg.AccessToken ||
		loaded.RefreshToken != cfg.RefreshToken ||
		loaded.ExpiresAt != cfg.ExpiresAt ||
		loaded.RefreshExpAt != cfg.RefreshExpAt ||
		loaded.Endpoint != cfg.Endpoint ||
		loaded.ClientID != cfg.ClientID ||
		loaded.ClientSecret != cfg.ClientSecret ||
		loaded.Username != cfg.Username {
		t.Errorf("loaded config %+v does not match saved %+v", loaded, cfg)
	}
}
