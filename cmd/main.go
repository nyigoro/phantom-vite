// phantom-vite/cmd/main.go

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Engine   string   `json:"engine"`
	Headless bool     `json:"headless"`
	Plugins  []string `json:"plugins"`
	Timeout  int      `json:"timeout"`
	Viewport struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"viewport"`
	Engines map[string]interface{} `json:"engines"`
}

type EngineStatus struct {
	Name      string
	Available bool
	Path      string
	Error     string
}

func loadConfig() Config {
	data, err := os.ReadFile("phantomvite.config.json")
	if err != nil {
		return Config{
			Engine:   "puppeteer",
			Headless: true,
			Timeout:  30000,
			Viewport: struct {
				Width  int `json:"width"`
				Height int `json:"height"`
			}{
				Width:  1920,
				Height: 1080,
			},
		}
	}
	var cfg Config
	json.Unmarshal(data, &cfg)

	if cfg.Engine == "" {
		cfg.Engine = "puppeteer"
	}
	if cfg.Viewport.Width == 0 {
		cfg.Viewport.Width = 1920
	}
	if cfg.Viewport.Height == 0 {
		cfg.Viewport.Height = 1080
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30000
	}

	return cfg
}

func checkEngineStatus() []EngineStatus {
	engines := []EngineStatus{
		{Name: "puppeteer", Available: false},
		{Name: "playwright", Available: false},
		{Name: "selenium", Available: false},
		{Name: "gemini", Available: false},
	}

	if puppeteerInstalled() {
		engines[0].Available = true
		engines[0].Path = "runtime/node_modules/puppeteer"
	} else {
		engines[0].Error = "Not installed. Run: cd runtime && npm install puppeteer"
	}

	if playwrightInstalled() {
		engines[1].Available = true
		engines[1].Path = "runtime/node_modules/playwright"
	} else {
		engines[1].Error = "Not installed. Run: cd runtime && npm install playwright"
	}

	if seleniumInstalled() {
		engines[2].Available = true
		engines[2].Path = "runtime-python"
	} else {
		engines[2].Error = "Not installed. Run: cd runtime-python && pip install selenium"
	}

	if _, err := exec.LookPath("gemini"); err == nil {
		engines[3].Available = true
		engines[3].Path = "system"
	} else {
		engines[3].Error = "Not installed. Run: npm install -g @google/gemini-cli"
	}

	return engines
}

// [main.go unchanged below this point — the rest of your CLI code remains intact]
