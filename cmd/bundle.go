package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
