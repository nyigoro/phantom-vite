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

func LoadPlugins(cfg Config) ([]string, error) {
	var loaded []string
	for _, path := range cfg.Plugins {
		abs, err := filepath.Abs(path)
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

func ExecutePluginHooks(hookName string, pluginPaths []string) {
	for _, plugin := range pluginPaths {
importPath := plugin
if os.PathSeparator == '\\' && strings.HasPrefix(plugin, "D:") {
	importPath = "file:///" + strings.ReplaceAll(plugin, "\\", "/")
}

cmd := exec.Command("node", "-e", fmt.Sprintf(`
  (async () => {
    try {
      const plugin = await import("%s");
      if (plugin.%s) await plugin.%s();
    } catch (e) {
      console.error("[Plugin Error]", e);
    }
  })()
`, importPath, hookName, hookName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = "runtime"
		_ = cmd.Run()
	}
}

func injectPluginContext() {
	cfg := loadConfig()
	if len(cfg.Plugins) > 0 {
		pluginEnv := strings.Join(cfg.Plugins, string(os.PathListSeparator))
		os.Setenv("PHANTOM_PLUGINS", pluginEnv)
	}
}

func init() {
	cfg := loadConfig()
	if plugins, err := LoadPlugins(cfg); err == nil {
		ExecutePluginHooks("onStart", plugins)
	}
	injectPluginContext()
}

func ExecutePluginHooksWithContext(hookName string, pluginPaths []string, context map[string]interface{}) {
	for _, plugin := range pluginPaths {
		serialized, _ := json.Marshal(context)
		cmd := exec.Command("node", "-e", fmt.Sprintf(`
			(async () => {
			  try {
				const plugin = await import("%s");
				if (plugin.%s) await plugin.%s(%s);
			  } catch (e) {
				console.error("[Plugin Error]", e);
			  }
			})()
		`, plugin, hookName, hookName, string(serialized)))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = "runtime"
		_ = cmd.Run()
	}
}

func writeTempScript(url string, engine string) (string, error) {
	var code string
	
	switch engine {
	case "puppeteer":
		code = fmt.Sprintf(`import puppeteer from 'puppeteer';
(async () => {
  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();
  await page.goto('%s');
  const title = await page.title();
  console.log("[Puppeteer] Page title:", title);
  await page.screenshot({ path: 'screenshot.png' });
  await browser.close();
})();`, url)
	case "playwright":
		code = fmt.Sprintf(`import { chromium } from 'playwright';
(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  await page.goto('%s');
  const title = await page.title();
  console.log("[Playwright] Page title:", title);
  await page.screenshot({ path: 'screenshot.png' });
  await browser.close();
})();`, url)
	default:
		return "", fmt.Errorf("unsupported engine for script generation: %s", engine)
	}

	tmpFile := fmt.Sprintf("phantom-open-%s.js", engine)
	err := os.WriteFile(tmpFile, []byte(code), 0644)
	return tmpFile, err
}

func runPageWithPlugins(script string, hooks []string) error {
	cfg := loadConfig()
	pluginPaths, _ := LoadPlugins(cfg)

	for _, hook := range hooks {
		ExecutePluginHooks(hook, pluginPaths)
	}

	err := runNodeScript(script)
	ExecutePluginHooks("onExit", pluginPaths)
	return err
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
		cmd := exec.Command(resolveCommand("python3"), "python/agent.py", prompt)
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
    "./plugins/performance.js"
  ]
}`)
			return
		}
		
		fmt.Println("📦 Plugin Status:")
		for _, path := range cfg.Plugins {
			if fileExists(path) {
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
			
			// Determine bundled file location
			baseName := strings.TrimSuffix(filepath.Base(script), ".ts")
			bundledScript := filepath.Join("dist", baseName+".js")
			
			if !fileExists(bundledScript) {
				fmt.Printf("❌ Bundled file not found: %s\n", bundledScript)
				
				// List dist directory contents for debugging
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
		
		if err := runNodeScript(script); err != nil {
			fmt.Printf("❌ Script execution failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ Script completed in %v\n", time.Since(start))
	}
}
