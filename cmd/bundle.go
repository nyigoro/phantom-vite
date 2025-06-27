package main

func bundleEngineScript(file string, engine string) error {
    // Auto-detect engine if empty
    if engine == "" {
        ext := strings.ToLower(filepath.Ext(file))
        switch ext {
        case ".ts", ".js":
            engine = "puppeteer" // default to Puppeteer for TypeScript/JavaScript
        case ".py":
            engine = "selenium"
        case ".gemini":
            engine = "gemini"
        default:
            return fmt.Errorf("could not auto-detect engine from file extension: %s", ext)
        }
        fmt.Println("🔍 Auto-detected engine:", engine)
    }

    // Proceed with engine logic
    switch engine {
    case "puppeteer", "playwright":
        cmd := exec.Command("npx", "vite", "build", "--config", "vite.config.js")
        cmd.Env = append(os.Environ(), "PHANTOM_ENTRY="+file)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        return cmd.Run()

    case "selenium", "gemini":
        out := filepath.Join("dist", filepath.Base(file))
        input, err := os.ReadFile(file)
        if err != nil {
            return err
        }
        return os.WriteFile(out, input, 0644)

    default:
        return fmt.Errorf("unsupported engine: %s", engine)
    }
}
