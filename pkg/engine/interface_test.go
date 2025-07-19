package engine

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Engine != "puppeteer" {
		t.Errorf("Expected default engine to be 'puppeteer', got %s", cfg.Engine)
	}

	if cfg.Headless != true {
		t.Errorf("Expected default headless to be true, got %t", cfg.Headless)
	}

	if cfg.Viewport.Width != 1920 {
		t.Errorf("Expected default viewport width to be 1920, got %d", cfg.Viewport.Width)
	}

	if cfg.Viewport.Height != 1080 {
		t.Errorf("Expected default viewport height to be 1080, got %d", cfg.Viewport.Height)
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30 seconds, got %s", cfg.Timeout)
	}

	if len(cfg.Plugins) != 0 {
		t.Errorf("Expected default plugins to be empty, got %d plugins", len(cfg.Plugins))
	}
}
