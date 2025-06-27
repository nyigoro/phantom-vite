package main

import "testing"

func writeTempScript(url string) (string, error) {
	cfg := loadConfig()
	pluginsJSON, _ := json.Marshal(cfg.Plugins)

	code := fmt.Sprintf(`import puppeteer from 'puppeteer';
import path from 'path';
import { pathToFileURL } from 'url';
import fs from 'fs';

console.log("[Phantom Vite] Loading plugins...");

const pluginPaths = %s;
const plugins = [];

for (const p of pluginPaths) {
  try {
    const module = await import(pathToFileURL(p).href);
    plugins.push(module);
  } catch (err) {
    console.error("[Plugin Load Error]", err.message);
  }
}

for (const plugin of plugins) {
  if (plugin.onStart) await plugin.onStart();
}

const browser = await puppeteer.launch({ headless: true });
const page = await browser.newPage();
await page.goto('%s');

for (const plugin of plugins) {
  if (plugin.onPageLoad) await plugin.onPageLoad(page);
}

await page.screenshot({ path: 'screenshot.png' });

await browser.close();

for (const plugin of plugins) {
  if (plugin.onExit) await plugin.onExit();
}
`, string(pluginsJSON), url)

	tmpFile := "phantom-open.js"
	err := os.WriteFile(tmpFile, []byte(code), 0644)
	return tmpFile, err
}

