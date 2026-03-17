package tests

import (
	"testing"

	"github.com/goclaw/goclaw/config"
)

func TestConfigLoad(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.App.Name != "GoClaw" {
		t.Errorf("Expected app name 'GoClaw', got '%s'", cfg.App.Name)
	}

	if cfg.App.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.App.Port)
	}
}
