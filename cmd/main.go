package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runNodeScript(script string, args ...string) error {
	cmdArgs := append([]string{script}, args...)
	cmd := exec.Command("node", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

func resolveCommand(name string) string {
	if os.PathSeparator == '\\' { // Windows
		if name == "python3" {
			return "python"
		}
	}
	return name
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
			fmt.Println("Usage: phantom-vite open <url>")
			os.Exit(1)
		}
		url := os.Args[2]
		fmt.Println("[Phantom Vite] Opening URL:", url)

		tmpScript, err := writeTempScript(url)
		if err != nil {
			fmt.Println("Error writing temp script:", err)
			os.Exit(1)
		}
		defer os.Remove(tmpScript)

		if err := runNodeScript(tmpScript); err != nil {
			fmt.Println("Error running Puppeteer:", err)
			os.Exit(1)
		}

	case "build":
		if err := runViteBuild(); err != nil {
			fmt.Println("Vite build failed:", err)
			os.Exit(1)
		}
        case "agent":
	       if len(os.Args) < 3 {
		      fmt.Println("Usage: phantom-vite agent <prompt>")
		       os.Exit(1)
	        }
	               prompt := strings.Join(os.Args[2:], " ")
	               cmd := exec.Command("python3", "python/agent.py", prompt)
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
	cmd := exec.Command("gemini", prompt)
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
        type Config struct {
	Headless bool `json:"headless"`
}

func loadConfig() Config {
	data, err := os.ReadFile("phantomvite.config.json")
	if err != nil {
		return Config{Headless: true} // default
	}
	var cfg Config
	json.Unmarshal(data, &cfg)
	return cfg
}

	default:
		script := os.Args[1]
		fmt.Println("[Phantom Vite] Running script:", script)
		if err := runNodeScript(script); err != nil {
			fmt.Println("Script error:", err)
			os.Exit(1)
		}
	}
}
