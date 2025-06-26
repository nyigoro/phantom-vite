package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: phantom-vite <script.js | command> [args]")
		os.Exit(1)
	}

	arg := os.Args[1]

	if isScript(arg) {
		runScript(arg)
	} else {
		fmt.Printf("[Phantom Vite] Command: %s\n", arg)
		switch arg {
		case "open":
			if len(os.Args) < 3 {
				fmt.Println("Usage: phantom-vite open <url>")
				return
			}
			url := os.Args[2]
			fmt.Printf("Opening URL: %s\n", url)
			// TODO: Add open logic
		default:
			fmt.Println("Unknown command:", arg)
		}
	}
}

func isScript(filename string) bool {
	return len(filename) > 3 && (filename[len(filename)-3:] == ".js" || filename[len(filename)-3:] == ".ts")
}

func runScript(scriptPath string) {
	fmt.Printf("[Phantom Vite] Running script: %s\n", scriptPath)

	// Adjust to full path if needed
	fullScriptPath, _ := filepath.Abs(scriptPath)
	fullWrapperPath, _ := filepath.Abs("engines/puppeteer-wrapper.js")

	cmd := exec.Command("node", fullWrapperPath, fullScriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("[Phantom Vite] Error running script:", err)
	}
}
