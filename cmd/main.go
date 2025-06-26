package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: phantom-vite <script.js>")
		os.Exit(1)
	}

	script := os.Args[1]

	fmt.Println("[Phantom Vite] Running script:", script)

	cmd := exec.Command("node", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running script:", err)
		os.Exit(1)
	}
}
