package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
)

supportedEngines := map[string]bool{
    "puppeteer": true,
    "playwright": true,
    "selenium": true,
    "gemini": true,
}

engine := "puppeteer" // default

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

func runScript(engine, file string) {
    var cmd *exec.Cmd

    switch engine {
    case "puppeteer", "playwright":
        cmd = exec.Command("node", file)
    case "selenium":
        cmd = exec.Command(resolvePythonCmd(), file)
    case "gemini":
        cmd = exec.Command("gemini", file)
    default:
        fmt.Println("Unsupported engine:", engine)
        return
    }

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err := cmd.Run()
    if err != nil {
        fmt.Println("Error running script:", err)
    }
}

func runNodeScript(script string) error {
	cmd := exec.Command("node", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ() // Pass PHANTOM_PLUGINS etc.
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
	cmd := exec.Command("python3", script) // Or python for Windows
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func validatePlugins(plugins []string) {
	fmt.Println("[Phantom Vite] Validating configured plugins...")
	for _, path := range plugins {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("❌ Plugin not found: %s\n", path)
		} else {
			fmt.Printf("✅ Plugin OK: %s\n", path)
		}
	}
}

func writeTempScript(url string) (string, error) {
	code := fmt.Sprintf(`import puppeteer from 'puppeteer';
(async () => {
  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();
  await page.goto('%s');
  const title = await page.title();
  console.log("[Phantom Vite] Title:", title);
  await page.screenshot({ path: 'screenshot.png' });
  await browser.close();
})();`, url)

	tmpFile := "phantom-open.js"
	err := os.WriteFile(tmpFile, []byte(code), 0644)
	return tmpFile, err
}

func runViteBuild() error {
	fmt.Println("[Phantom Vite] Running Vite build...")
	cmd := exec.Command("npx", "vite", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runViteBundle(entry string) error {
	base := strings.TrimSuffix(filepath.Base(entry), filepath.Ext(entry))
	outFile := fmt.Sprintf("dist/%s.js", base)

	fmt.Println("[Phantom Vite] Bundling:", entry)
	cmd := exec.Command("npx", "vite", "build",
		"--config", "vite.config.js",
		"--entry", entry,
		"--outDir", "dist",
		"--rollupOptions.output.file", outFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func puppeteerInstalled() bool {
	info, err := os.Stat("node_modules/puppeteer")
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

// ✅ Config definition moved outside main
type Config struct {
	Headless bool `json:"headless"`
        Plugins  []string `json:"plugins"`
}

func loadConfig() Config {
	data, err := os.ReadFile("phantomvite.config.json")
	if err != nil {
		return Config{Headless: true}
	}
	var cfg Config
	json.Unmarshal(data, &cfg)
	return cfg
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  phantom-vite open <url>")
		fmt.Println("  phantom-vite build")
		fmt.Println("  phantom-vite <script.js>")
		os.Exit(1)
	}

	switch os.Args[1] {
case "open":
    if len(os.Args) < 3 {
        fmt.Println("Usage: phantom-vite open <url> [--engine engine]")
        return
    }

    url := os.Args[2]
    engine := "puppeteer" // default

    for i := 3; i < len(os.Args); i++ {
        if os.Args[i] == "--engine" && i+1 < len(os.Args) {
            engine = os.Args[i+1]
        }
    }

    // Load plugins
    cfg := loadConfig()
    pluginEnv := strings.Join(cfg.Plugins, ",")

    // Generate script content
    scriptContent, err := GenerateTempScript(engine, url, pluginEnv)
    if err != nil {
        fmt.Println("Error generating script:", err)
        return
    }

    tempFile, err := writeTempFile(scriptContent, engine)
    if err != nil {
        fmt.Println("Error writing temp script:", err)
        return
    }

    defer os.Remove(tempFile)

    fmt.Println("[Phantom Vite] Running script:", tempFile)
    runScript(engine, tempFile)
        case "engines":
    fmt.Println("🔧 Supported engines:")
    fmt.Println(" - puppeteer  (Node.js, full Chrome control via DevTools protocol)")
    fmt.Println(" - playwright (Node.js, cross-browser automation)")
    fmt.Println(" - selenium   (Python, WebDriver-based, cross-language)")
    fmt.Println(" - gemini     (CLI, Gemini tool integration for testing)")
    return
	case "build": {
		if err := runViteBuild(); err != nil {
			fmt.Println("Vite build failed:", err)
			os.Exit(1)
		}
	}
		case "bundle": {
	if len(os.Args) < 3 {
		fmt.Println("Usage: phantom-vite bundle <scripts/your-file.ts>")
		os.Exit(1)
	}
	entry := os.Args[2]
	if err := runViteBundle(entry); err != nil {
		fmt.Println("Vite bundle failed:", err)
		os.Exit(1)
	}
}
	case "agent": {
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
	}
	case "gemini": {
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
	}	
	case "serve": {
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
	}
        case "plugins": {
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
}
	default: {
		script := os.Args[1]
		fmt.Println("[Phantom Vite] Running script:", script)
		if err := runNodeScript(script); err != nil {
			fmt.Println("Script error:", err)
			os.Exit(1)
		}
	}
	}
}
