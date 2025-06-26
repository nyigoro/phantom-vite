// engines/puppeteer-wrapper.js
const fs = require("fs");
const path = require("path");

// Get script file path from CLI args
const scriptPath = process.argv[2];

if (!scriptPath || !fs.existsSync(scriptPath)) {
  console.error("[Phantom Vite] Script file not found:", scriptPath);
  process.exit(1);
}

// Load and run the bundled script
(async () => {
  console.log(`[Phantom Vite] Executing: ${scriptPath}`);
  const script = fs.readFileSync(scriptPath, "utf-8");
  eval(script); // TEMP: for dev only — will replace with better runner
})();

