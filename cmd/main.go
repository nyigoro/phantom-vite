package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: phantom-vite <script.js | command> [args]")
		os.Exit(1)
	}

	arg := os.Args[1]

	if isScript(arg) {
		fmt.Printf("[Phantom Vite] Running script: %s\n", arg)
		// TODO: call script runner
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
			// TODO: pass to engine
		default:
			fmt.Println("Unknown command:", arg)
		}
	}
}

func isScript(filename string) bool {
	return len(filename) > 3 && (filename[len(filename)-3:] == ".js" || filename[len(filename)-3:] == ".ts")
}
