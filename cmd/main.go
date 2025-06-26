package main

import (
	"fmt"
	"os"
	"os/exec"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: phantom-vite <script.js> or phantom-vite open <url>")
		os.Exit(1)
	}

	if os.Args[1] == "open" && len(os.Args) >= 3 {
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
		return
	}

	// Default mode: run JS script
	script := os.Args[1]
	fmt.Println("[Phantom Vite] Running script:", script)
	if err := runNodeScript(script); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
