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

type PluginConfig struct {
	Path    string                 `json:"path"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type Config struct {
	Engine   string         `json:"engine"`
	Timeout  int            `json:"timeout"` 
	Headless bool           `json:"headless"`
	Plugins  []PluginConfig `json:"plugins"`
	Entries  []string       `json:"entries"`
	Viewport struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"viewport"`
}

type EngineStatus struct {
	Name      string
	Available bool
	Path      string
	Error     string
}

type PluginContext struct {
	Engine   string   `json:"engine"`
	Headless bool     `json:"headless"`
	Plugins  []PluginConfig `json:"plugins"`
	Timeout  int      `json:"timeout"`
	Viewport struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"viewport"`
	Meta struct {
		Command string `json:"command"`
		Script  string `json:"script,omitempty"`
		URL     string `json:"url,omitempty"`
	} `json:"meta"`
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

func LoadPlugins(cfg Config) ([]string, error) {
	var loaded []string
	for _, path := range cfg.Plugins {
		abs, err := filepath.Abs(path.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve plugin path %s: %v", path, err)
		}
		if _, err := os.Stat(abs); os.IsNotExist(err) {
			return nil, fmt.Errorf("plugin not found: %s", abs)
		}
		loaded = append(loaded, abs)
	}
	return loaded, nil
}

func injectPluginContext() {
	cfg := loadConfig()
	if len(cfg.Plugins) > 0 {
		var pluginPaths []string
		for _, plugin := range cfg.Plugins {
			pluginPaths = append(pluginPaths, plugin.Path)
		}
		pluginEnv := strings.Join(pluginPaths, string(os.PathListSeparator))
		os.Setenv("PHANTOM_PLUGINS", pluginEnv)
	}
}

func ExecutePluginHooksWithContext(hookName string, pluginPaths []string, context PluginContext) {
	for _, plugin := range pluginPaths {
		// Normalize for Windows paths with drive letters (D:\ etc)
		importPath := plugin
		if os.PathSeparator == '\\' && strings.HasPrefix(plugin, "D:") {
			importPath = "file:///" + strings.ReplaceAll(plugin, "\\", "/")
		}

		serialized, _ := json.Marshal(context)

		cmd := exec.Command("node", "-e", fmt.Sprintf(`(async () => {
		  try {
			const plugin = await import("%s");
			if (plugin.%s) await plugin.%s(%s);
		  } catch (e) {
			console.error("[Plugin Error]", e);
		  }
		})()`, importPath, hookName, hookName, string(serialized)))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = "runtime"
		_ = cmd.Run()
	}
}

func runPageWithPlugins(script string, hooks []string, command string) error {
	cfg := loadConfig()
	pluginPaths, _ := LoadPlugins(cfg)

	ctx := PluginContext{
	// existing fields
	Engine:   cfg.Engine,
	Headless: cfg.Headless,
	Plugins:  cfg.Plugins,
	Timeout:  cfg.Timeout,
	Viewport: cfg.Viewport,
	Meta: struct {
		Command string `json:"command"`
		Script  string `json:"script,omitempty"`
		URL     string `json:"url,omitempty"`
	}{
		Command: "run",
		Script:  script,
	},
}

	for _, hook := range hooks {
		ExecutePluginHooksWithContext(hook, pluginPaths, ctx)
	}

	err := runNodeScript(script)
	ExecutePluginHooksWithContext("onExit", pluginPaths, ctx)
	return err
}

// Add this function to your main.go file
func writeTempScript(url string, engine string) (string, error) {
	var code string
	
	switch engine {
	case "puppeteer":
		code = fmt.Sprintf(`import puppeteer from 'puppeteer';

const pluginPaths = process.env["PHANTOM_PLUGINS"]?.split(",") ?? [];
const plugins = [];

for (const path of pluginPaths) {
  try {
    const mod = await import(path);
    plugins.push(mod);
  } catch (e) {
    console.error("[Phantom Vite] Failed to load plugin:", path, e);
  }
}

(async () => {
  for (const p of plugins) {
    if (typeof p.onStart === 'function') await p.onStart();
  }

  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();

  const url = '%s';
  await page.goto(url);

  for (const p of plugins) {
    if (typeof p.onPageLoad === 'function') await p.onPageLoad(page);
  }

  const title = await page.title();
  console.log("[Phantom Vite] Title:", title);
  await page.screenshot({ path: 'screenshot.png' });

  await browser.close();

  for (const p of plugins) {
    if (typeof p.onExit === 'function') await p.onExit();
  }
})();`, url)

	case "playwright":
		code = fmt.Sprintf(`import { chromium } from 'playwright';

const pluginPaths = process.env["PHANTOM_PLUGINS"]?.split(",") ?? [];
const plugins = [];

for (const path of pluginPaths) {
  try {
    const mod = await import(path);
    plugins.push(mod);
  } catch (e) {
    console.error("[Phantom Vite] Failed to load plugin:", path, e);
  }
}

(async () => {
  for (const p of plugins) {
    if (typeof p.onStart === 'function') await p.onStart();
  }

  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  const url = '%s';
  await page.goto(url);

  for (const p of plugins) {
    if (typeof p.onPageLoad === 'function') await p.onPageLoad(page);
  }

  const title = await page.title();
  console.log("[Phantom Vite] Title:", title);
  await page.screenshot({ path: 'screenshot.png' });

  await browser.close();

  for (const p of plugins) {
    if (typeof p.onExit === 'function') await p.onExit();
  }
})();`, url)

	case "selenium":
		code = fmt.Sprintf(`from selenium import webdriver
from selenium.webdriver.chrome.options import Options
import os
import json

# Load plugins (simplified for Python)
plugin_paths = os.environ.get("PHANTOM_PLUGINS", "").split(",") if os.environ.get("PHANTOM_PLUGINS") else []

# Setup Chrome options
chrome_options = Options()
chrome_options.add_argument("--headless")
chrome_options.add_argument("--no-sandbox")
chrome_options.add_argument("--disable-dev-shm-usage")

try:
    driver = webdriver.Chrome(options=chrome_options)
    
    url = "%s"
    driver.get(url)
    
    title = driver.title
    print(f"[Phantom Vite] Title: {title}")
    
    driver.save_screenshot("screenshot.png")
    
finally:
    driver.quit()
`, url)

	default:
		return "", fmt.Errorf("unsupported engine: %s", engine)
	}

	var tmpFile string
	if engine == "selenium" {
		tmpFile = "phantom-open.py"
	} else {
		tmpFile = "phantom-open.js"
	}
	
	err := os.WriteFile(tmpFile, []byte(code), 0644)
	return tmpFile, err
}

func validateEngine(engine string) error {
	status := checkEngineStatus()
	for _, s := range status {
		if s.Name == engine {
			if !s.Available {
				return fmt.Errorf("❌ Engine '%s' is not available: %s", engine, s.Error)
			}
			return nil
		}
	}
	return fmt.Errorf("❌ Unknown engine: %s", engine)
}

func runEngineScript(path, engine string) error {
	if engine == "selenium" {
		cmd := exec.Command(resolveCommand("python3"), path)
		cmd.Dir = "runtime-python"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	cmd := exec.Command("node", path)
	cmd.Dir = "runtime"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runNodeScript(script string) error {
	cmd := exec.Command("node", script)

	// Handle absolute path
	if filepath.IsAbs(script) {
		cmd.Dir = filepath.Dir(script)
		script = filepath.Base(script)
		cmd.Args = []string{"node", script}
	} else if strings.HasPrefix(script, "dist/") || strings.HasPrefix(script, "dist\\") {
		// Bundled output is in project root/dist
		cmd.Dir = "."
	} else {
		// Relative custom scripts assumed in runtime
		cmd.Dir = "runtime"
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeContextFile(cfg Config, url string) (string, error) {
	context := map[string]interface{}{
		"engine":   cfg.Engine,
		"headless": cfg.Headless,
		"viewport": cfg.Viewport,
		"timeout":  cfg.Timeout,
		"plugins":  cfg.Plugins,
		"url":      url,
	}

	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return "", err
	}

	path := "phantom.context.json"
	err = os.WriteFile(path, data, 0644)
	return path, err
}

func runViteBuild() error {
	fmt.Println("🔧 [Phantom Vite] Running Vite build...")
	cmd := exec.Command("npx", "vite", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runViteBundle(entry string) error {
	fmt.Printf("📦 [Phantom Vite] Bundling: %s\n", entry)
	
	// Create a temporary Vite config for this specific entry
	tempConfig := fmt.Sprintf(`
import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    lib: {
      entry: '%s',
      name: 'PhantomViteScript',
      fileName: (format) => '%s.js'
    },
    rollupOptions: {
      external: ['puppeteer', 'playwright', 'selenium-webdriver'],
      output: {
        globals: {
          'puppeteer': 'puppeteer',
          'playwright': 'playwright',
          'selenium-webdriver': 'selenium'
        }
      }
    }
  }
});
`, entry, strings.TrimSuffix(filepath.Base(entry), filepath.Ext(entry)))
	
	configFile := "vite.config.temp.js"
	err := os.WriteFile(configFile, []byte(tempConfig), 0644)
	if err != nil {
		return fmt.Errorf("failed to create temp config: %w", err)
	}
	defer os.Remove(configFile)
	
	cmd := exec.Command("npx", "vite", "build", "--config", configFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func puppeteerInstalled() bool {
	info, err := os.Stat("runtime/node_modules/puppeteer")
	return err == nil && info.IsDir()
}

func playwrightInstalled() bool {
	info, err := os.Stat("runtime/node_modules/playwright")
	return err == nil && info.IsDir()
}

func seleniumInstalled() bool {
	info, err := os.Stat("runtime-python")
	return err == nil && info.IsDir()
}

func resolveCommand(name string) string {
	if os.PathSeparator == '\\' {
		if name == "python3" {
			return "python"
		}
	}
	return name
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func printUsage() {
	fmt.Println("🕴️  Phantom Vite - Headless Browser CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  phantom-vite open <url> [--engine <engine>]")
	fmt.Println("  phantom-vite build")
	fmt.Println("  phantom-vite bundle <file>")
	fmt.Println("  phantom-vite serve <file>")
	fmt.Println("  phantom-vite doctor")
	fmt.Println("  phantom-vite engines")
	fmt.Println("  phantom-vite agent <prompt>")
	fmt.Println("  phantom-vite gemini <prompt>")
	fmt.Println("  phantom-vite plugins")
	fmt.Println("  phantom-vite <script.js>")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  phantom-vite open https://example.com")
	fmt.Println("  phantom-vite open https://example.com --engine playwright")
	fmt.Println("  phantom-vite build")
	fmt.Println("  phantom-vite script.ts")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg := loadConfig()
	
	// Parse engine flag
	engine := cfg.Engine
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--engine" && i+1 < len(os.Args) {
			engine = os.Args[i+1]
			break
		}
	}

	switch os.Args[1] {
	case "doctor":
		fmt.Println("🏥 [Phantom Vite] Health Check")
		fmt.Println()
		
		// Check Go version
		if cmd := exec.Command("go", "version"); cmd.Run() == nil {
			fmt.Println("✅ Go runtime available")
		} else {
			fmt.Println("❌ Go runtime not found")
		}
		
		// Check Node.js version
		if cmd := exec.Command("node", "--version"); cmd.Run() == nil {
			fmt.Println("✅ Node.js runtime available")
		} else {
			fmt.Println("❌ Node.js runtime not found")
		}
		
		// Check Python version
		if cmd := exec.Command(resolveCommand("python3"), "--version"); cmd.Run() == nil {
			fmt.Println("✅ Python runtime available")
		} else {
			fmt.Println("❌ Python runtime not found")
		}
		
		// Check engines
		fmt.Println()
		fmt.Println("🔧 Engine Status:")
		statuses := checkEngineStatus()
		for _, status := range statuses {
			if status.Available {
				fmt.Printf("  ✅ %s (%s)\n", status.Name, status.Path)
			} else {
				fmt.Printf("  ❌ %s - %s\n", status.Name, status.Error)
			}
		}
		
		// Check configuration
		fmt.Println()
		fmt.Println("⚙️  Configuration:")
		fmt.Printf("  Default engine: %s\n", cfg.Engine)
		fmt.Printf("  Headless mode: %t\n", cfg.Headless)
		fmt.Printf("  Viewport: %dx%d\n", cfg.Viewport.Width, cfg.Viewport.Height)
		fmt.Printf("  Timeout: %dms\n", cfg.Timeout)
		
		return

	case "open":
		args := os.Args[2:]
		if len(args) < 1 {
			fmt.Println("Usage: phantom-vite open <url> [--engine <engine>]")
			return
		}
		
		url := args[0]
		
		// Validate engine
		if err := validateEngine(engine); err != nil {
			fmt.Println(err)
			fmt.Println("💡 Run 'phantom-vite doctor' to check your setup")
			return
		}
		
		scriptPath, err := writeTempScript(url, engine)
		if err != nil {
			fmt.Printf("❌ Failed to generate script: %v\n", err)
			return
		}
		defer os.Remove(scriptPath) // Clean up temp file
		
		fmt.Printf("🚀 Opening %s with %s engine...\n", url, engine)
		
		start := time.Now()
pluginPaths, _ := LoadPlugins(cfg)
ctx := PluginContext{
	Engine:   cfg.Engine,
	Headless: cfg.Headless,
	Plugins:  cfg.Plugins,
	Timeout:  cfg.Timeout,
	Viewport: cfg.Viewport,
	Meta: struct {
		Command string `json:"command"`
		Script  string `json:"script,omitempty"`
		URL     string `json:"url,omitempty"`
	}{
		Command: "open",
	        URL:     url,
	        Script:  scriptPath,
	},
}
ExecutePluginHooksWithContext("onStart", pluginPaths, ctx)
ctx.Meta.Command = "open"
ctx.Meta.URL = url
ExecutePluginHooksWithContext("onStart", pluginPaths, ctx)

contextPath, err := writeContextFile(cfg, url)
if err == nil {
	os.Setenv("PHANTOM_CONTEXT_PATH", contextPath)
}

		if err := runEngineScript(scriptPath, engine); err != nil {
			fmt.Printf("❌ Script execution failed: %v\n", err)
			return
		}

		defer os.Remove("phantom.context.json")
		fmt.Printf("✅ Completed in %v\n", time.Since(start))

	case "engines":
		fmt.Println("🔧 Supported Engines:")
		statuses := checkEngineStatus()
		for _, status := range statuses {
			availability := "❌"
			if status.Available {
				availability = "✅"
			}
			fmt.Printf("  %s %s", availability, status.Name)
			
			switch status.Name {
			case "puppeteer":
				fmt.Print(" - Node.js, full Chrome control via DevTools protocol")
			case "playwright":
				fmt.Print(" - Node.js, cross-browser automation (Chrome, Firefox, Safari)")
			case "selenium":
				fmt.Print(" - Python, WebDriver-based, cross-language support")
			case "gemini":
				fmt.Print(" - CLI, Google AI integration for intelligent automation")
			}
			
			if !status.Available {
				fmt.Printf("\n    💡 %s", status.Error)
			}
			fmt.Println()
		}
		return

	case "build":
		fmt.Println("🔧 [Phantom Vite] Building project...")
		start := time.Now()
		if err := runViteBuild(); err != nil {
			fmt.Printf("❌ Build failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Build completed in %v\n", time.Since(start))

	case "bundle":
		args := os.Args[2:]
		if len(args) < 1 {
			fmt.Println("Usage: phantom-vite bundle <file> [--engine <engine>]")
			return
		}
		
		inputFile := args[0]
		if !fileExists(inputFile) {
			fmt.Printf("❌ File not found: %s\n", inputFile)
			return
		}
		
		fmt.Printf("📦 Bundling %s for %s engine...\n", inputFile, engine)
		
		start := time.Now()
		if err := runViteBundle(inputFile); err != nil {
			fmt.Printf("❌ Bundling failed: %v\n", err)
			return
		}
		
		fmt.Printf("✅ Bundling completed in %v\n", time.Since(start))

	case "serve":
		if len(os.Args) < 3 {
			fmt.Println("Usage: phantom-vite serve <file>")
			os.Exit(1)
		}
		file := os.Args[2]
		
		if !fileExists(file) {
			fmt.Printf("❌ File not found: %s\n", file)
			return
		}
		
		fmt.Printf("🌐 Serving %s...\n", file)
		cmd := exec.Command("npx", "vite", "preview", "--config", file)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Serve error: %v\n", err)
			os.Exit(1)
		}

case "agent":
	if len(os.Args) < 3 {
		fmt.Println("Usage: phantom-vite agent <prompt>")
		os.Exit(1)
	}
	prompt := strings.Join(os.Args[2:], " ")
	fmt.Printf("🤖 Launching AI agent with prompt: %s\n", prompt)

	pluginPaths, _ := LoadPlugins(cfg)
	context := PluginContext{
	Engine:   cfg.Engine,
	Headless: cfg.Headless,
	Plugins:  cfg.Plugins,
	Timeout:  cfg.Timeout,
	Viewport: cfg.Viewport,
	Meta: struct {
		Command string `json:"command"`
		Script  string `json:"script,omitempty"`
		URL     string `json:"url,omitempty"`
	}{
		Command: "agent",
	},
}
	context.Meta.Command = "agent"
ExecutePluginHooksWithContext("onStart", pluginPaths, context)

	cmd := exec.Command(resolveCommand("python3"), "python/agent.py", prompt)
	cmd.Dir = filepath.Join(filepath.Dir(os.Args[0]), "..") // Set working directory to project root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Agent error: %v\n", err)
		os.Exit(1)
	}

case "gemini":
	if len(os.Args) < 3 {
		fmt.Println("Usage: phantom-vite gemini <prompt>")
		os.Exit(1)
	}
	prompt := strings.Join(os.Args[2:], " ")
	fmt.Printf("✨ Passing to Gemini CLI: %s\n", prompt)

	pluginPaths, _ := LoadPlugins(cfg)
	context := PluginContext{
	Engine:   cfg.Engine,
	Headless: cfg.Headless,
	Plugins:  cfg.Plugins,
	Timeout:  cfg.Timeout,
	Viewport: cfg.Viewport,
	Meta: struct {
		Command string `json:"command"`
		Script  string `json:"script,omitempty"`
		URL     string `json:"url,omitempty"`
	}{
		Command: "gemini",
	},
}
	context.Meta.Command = "gemini"
ExecutePluginHooksWithContext("onStart", pluginPaths, context)

	cmd := exec.Command("gemini", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Gemini CLI error: %v\n", err)
	}

	case "plugins":
		cfg := loadConfig()
		if len(cfg.Plugins) == 0 {
			fmt.Println("📦 No plugins defined in phantomvite.config.json")
			fmt.Println("💡 Add plugins to your config file:")
			fmt.Println(`{
  "plugins": [
    "./plugins/seo.js",
    "./plugins/logger.js",
    "./plugins/performance.js"
  ]
}`)
			return
		}
		
		fmt.Println("📦 Plugin Status:")
		for _, path := range cfg.Plugins {
			if fileExists(path.Path) {
				fmt.Printf("  ✅ %s\n", path)
			} else {
				fmt.Printf("  ❌ %s (not found)\n", path)
			}
		}

default:
	script := os.Args[1]

	if !fileExists(script) {
		fmt.Printf("❌ Script not found: %s\n", script)
		fmt.Println("💡 Make sure the file path is correct")
		return
	}

	ext := filepath.Ext(script)

	if ext == ".ts" {
		fmt.Printf("🔧 TypeScript detected, bundling %s...\n", script)

		if err := runViteBundle(script); err != nil {
			fmt.Printf("❌ Failed to bundle: %v\n", err)
			return
		}

		baseName := strings.TrimSuffix(filepath.Base(script), ".ts")
		bundledScript := filepath.Join("dist", baseName+".js")

		if !fileExists(bundledScript) {
			fmt.Printf("❌ Bundled file not found: %s\n", bundledScript)
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

	fmt.Printf("🚀 Running script: %s\n", script)
	start := time.Now()

	// ✅ Inject context before running the script
	context := PluginContext{
	Engine:   cfg.Engine,
	Headless: cfg.Headless,
	Plugins:  cfg.Plugins,
	Timeout:  cfg.Timeout,
	Viewport: cfg.Viewport,
	Meta: struct {
		Command string `json:"command"`
		Script  string `json:"script,omitempty"`
		URL     string `json:"url,omitempty"`
	}{
		Command: "run",
		Script:  script,
	},
}
	pluginPaths, _ := LoadPlugins(cfg)
	ExecutePluginHooksWithContext("onStart", pluginPaths, context)

	if err := runNodeScript(script); err != nil {
		fmt.Printf("❌ Script execution failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Script completed in %v\n", time.Since(start))
	}
}
