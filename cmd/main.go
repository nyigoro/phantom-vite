package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Headless bool     `json:"headless"`
	Plugins  []string `json:"plugins"`
	Viewport struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"viewport"`
}

func loadConfig() Config {
	data, err := os.ReadFile("phantomvite.config.json")
	if err != nil {
		return Config{
			Headless: true,
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
	
	// Set default viewport if not specified
	if cfg.Viewport.Width == 0 {
		cfg.Viewport.Width = 1920
	}
	if cfg.Viewport.Height == 0 {
		cfg.Viewport.Height = 1080
	}
	
	return cfg
}

func writeTempScript(url string, engine string) (string, error) {
	code := fmt.Sprintf(`import puppeteer from 'puppeteer';
(async () => {
  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();
  await page.goto('%s');
  const title = await page.title();
  console.log("[Puppeteer] Page title:", title);
  await page.screenshot({ path: 'screenshot.png' });
  await browser.close();
})();`, url)

	tmpFile := "phantom-open.js"
	err := os.WriteFile(tmpFile, []byte(code), 0644)
	return tmpFile, err
}

func runEngineScript(path, engine string) {
	if engine == "selenium" {
		cmd := exec.Command(resolveCommand("python3"), path)
		cmd.Dir = "runtime-python"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	cmd := exec.Command("node", path)
	cmd.Dir = "runtime"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Updated runNodeScript to handle both regular scripts and bundled scripts
func runNodeScript(script string) error {
	cmd := exec.Command("node", script)
	
	// Only change to runtime directory for scripts that aren't bundled
	if !strings.HasPrefix(script, "dist/") && !strings.HasPrefix(script, "../dist/") {
		cmd.Dir = "runtime"
	}
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func runGeminiPrompt(prompt string) error {
	cmd := exec.Command("gemini", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func runPythonScript(script string) error {
	cmd := exec.Command("python3", script)
	cmd.Dir = "runtime-python"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func runViteBuild() error {
	fmt.Println("[Phantom Vite] Running Vite build...")
	cmd := exec.Command("npx", "vite", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runViteBundle(entry string) error {
	fmt.Println("[Phantom Vite] Bundling:", entry)
	cmd := exec.Command("npx", "vite", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func puppeteerInstalled() bool {
	info, err := os.Stat("runtime/node_modules/puppeteer")
	if err != nil {
		return false
	}
	return info.IsDir()
}

func resolveCommand(name string) string {
	if os.PathSeparator == '\\' {
		if name == "python3" {
			return "python"
		}
	}
	return name
}

func bundleIfTs(file string) error {
	ext := filepath.Ext(file)
	if ext == ".ts" {
		fmt.Println("[Phantom Vite] Detected TypeScript, bundling...")
		return runViteBundle(file)
	}
	return nil
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  phantom-vite open <url>")
		fmt.Println("  phantom-vite build")
		fmt.Println("  phantom-vite <script.js>")
		os.Exit(1)
	}

	supportedEngines := map[string]bool{
		"puppeteer":  true,
		"playwright": true,
		"selenium":   true,
		"gemini":     true,
	}

	engine := "puppeteer"
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--engine" && i+1 < len(os.Args) {
			engine = os.Args[i+1]
		}
	}

	if !supportedEngines[engine] {
		fmt.Println("❌ Unsupported engine:", engine)
		fmt.Println("✅ Supported engines: puppeteer, playwright, selenium, gemini")
		return
	}

	args := os.Args[2:]

	switch os.Args[1] {
	case "open":
		if len(args) < 1 {
			fmt.Println("Usage: phantom-vite open <url> [--engine engineName]")
			return
		}
		url := args[0]
		if len(args) > 2 && args[1] == "--engine" {
			engine = args[2]
		}
		scriptPath, err := writeTempScript(url, engine)
		if err != nil {
			fmt.Println("❌ Failed to generate script:", err)
			return
		}
		fmt.Println("🚀 Running script:", scriptPath)
		runEngineScript(scriptPath, engine)

	case "engines":
		fmt.Println("🔧 Supported engines:")
		fmt.Println(" - puppeteer  (Node.js, full Chrome control via DevTools protocol)")
		fmt.Println(" - playwright (Node.js, cross-browser automation)")
		fmt.Println(" - selenium   (Python, WebDriver-based, cross-language)")
		fmt.Println(" - gemini     (CLI, Gemini tool integration for testing)")
		return

	case "build":
		if err := runViteBuild(); err != nil {
			fmt.Println("Vite build failed:", err)
			os.Exit(1)
		}

	case "bundle":
		if len(args) < 1 {
			fmt.Println("Usage: phantom-vite bundle <file> [--engine engineName]")
			return
		}
		inputFile := args[0]
		if len(args) > 2 && args[1] == "--engine" {
			engine = args[2]
		}
		err := runViteBundle(inputFile)
		if err != nil {
			fmt.Println("❌ Bundling failed:", err)
			return
		}
		fmt.Println("📦 Bundled", inputFile, "for engine:", engine)

	case "agent":
		if len(os.Args) < 3 {
			fmt.Println("Usage: phantom-vite agent <prompt>")
			os.Exit(1)
		}
		prompt := strings.Join(os.Args[2:], " ")
		cmd := exec.Command(resolveCommand("python3"), "python/agent.py", prompt)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		fmt.Println("[Phantom Vite] Launching AI agent...")
		if err := cmd.Run(); err != nil {
			fmt.Println("Agent error:", err)
			os.Exit(1)
		}

	case "gemini":
		if len(os.Args) < 3 {
			fmt.Println("Usage: phantom-vite gemini <prompt>")
			os.Exit(1)
		}
		prompt := strings.Join(os.Args[2:], " ")
		cmd := exec.Command(resolveCommand("gemini"), prompt)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		fmt.Println("[Phantom Vite] Passing to Gemini CLI...")
		if err := cmd.Run(); err != nil {
			fmt.Println("Gemini CLI error:", err)
		}

	case "serve":
		if len(os.Args) < 3 {
			fmt.Println("Usage: phantom-vite serve <file>")
			os.Exit(1)
		}
		file := os.Args[2]
		cmd := exec.Command("npx", "vite", "preview", "--config", file)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		fmt.Println("[Phantom Vite] Serving file:", file)
		if err := cmd.Run(); err != nil {
			fmt.Println("Serve error:", err)
			os.Exit(1)
		}

	case "plugins":
		cfg := loadConfig()
		if len(cfg.Plugins) == 0 {
			fmt.Println("[Phantom Vite] No plugins defined in phantomvite.config.json")
			return
		}
		fmt.Println("[Phantom Vite] Plugin Check:")
		for _, path := range cfg.Plugins {
			if _, err := os.Stat(path); err != nil {
				fmt.Printf("  ❌ %s (not found)\n", path)
			} else {
				fmt.Printf("  ✅ %s\n", path)
			}
		}

	default:
		script := os.Args[1]
		ext := filepath.Ext(script)
		
		if ext == ".ts" {
			fmt.Println("[Phantom Vite] Detected TypeScript file, bundling...")
			if err := bundleIfTs(script); err != nil {
				fmt.Println("❌ Failed to bundle:", err)
				return
			}
			
			// Determine the expected output file name
			baseName := strings.TrimSuffix(filepath.Base(script), ".ts")
			bundledScript := "dist/" + baseName + ".js"
			
			// Check if the bundled file exists
			if !fileExists(bundledScript) {
				fmt.Printf("❌ Bundled file not found: %s\n", bundledScript)
				fmt.Println("💡 Make sure Vite build completed successfully")
				
				// Try to list what's actually in the dist directory
				if files, err := os.ReadDir("dist"); err == nil {
					fmt.Println("📁 Files in dist directory:")
					for _, file := range files {
						fmt.Printf("  - %s\n", file.Name())
					}
				}
				return
			}
			
			script = bundledScript
			fmt.Printf("✅ Using bundled script: %s\n", script)
		}
		
		fmt.Println("[Phantom Vite] Running script:", script)
		if err := runNodeScript(script); err != nil {
			fmt.Println("Script error:", err)
			os.Exit(1)
		}
	}
}
