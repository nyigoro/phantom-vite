package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPlugins_ValidPaths(t *testing.T) {
	tempPlugin := "plugins/test-plugin.js"

	// Create dummy plugin file
	err := os.MkdirAll(filepath.Dir(tempPlugin), 0755)
	if err != nil {
		t.Fatalf("failed to create plugins dir: %v", err)
	}
	err = os.WriteFile(tempPlugin, []byte("export const onStart = () => console.log('loaded');"), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy plugin: %v", err)
	}
	defer os.Remove(tempPlugin)

	cfg := Config{Plugins: []string{tempPlugin}}

	paths, err := LoadPlugins(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(paths) != 1 {
		t.Fatalf("expected 1 plugin path, got %d", len(paths))
	}
}

func TestLoadPlugins_InvalidPath(t *testing.T) {
	cfg := Config{Plugins: []string{"nonexistent/plugin.js"}}

	_, err := LoadPlugins(cfg)
	if err == nil {
		t.Fatal("expected error for missing plugin, got nil")
	}
}
