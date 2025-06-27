func bundleEngineScript(file, engine string) error {
    switch engine {
    case "puppeteer", "playwright":
        cmd := exec.Command("npx", "vite", "build", "--config", "vite.config.js")
        cmd.Env = append(os.Environ(), "PHANTOM_ENTRY="+file)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        return cmd.Run()

    case "selenium":
        // Just copy to dist/ for now; Python scripts don't bundle
        out := filepath.Join("dist", filepath.Base(file))
        input, err := os.ReadFile(file)
        if err != nil {
            return err
        }
        return os.WriteFile(out, input, 0644)

    case "gemini":
        // Gemini scripts are plain text — copy to dist/
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
