// engine_scripts.go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func writePuppeteerScript(url string) (string, error) {
	cfg := loadConfig()

	script := fmt.Sprintf(`import puppeteer from 'puppeteer';

const pluginPaths = process.env["PHANTOM_PLUGINS"]?.split(',') ?? [];
const plugins = [];

for (const path of pluginPaths) {
  try {
    const mod = await import(path);
    plugins.push(mod);
  } catch (e) {
    console.error("[Phantom Vite] Failed to load plugin:", path, e);
  }
}

(async () => {
  for (const p of plugins) {
    if (typeof p.onStart === 'function') await p.onStart();
  }

  const browser = await puppeteer.launch({ headless: %v });
  const page = await browser.newPage();
  await page.setViewport({ width: %d, height: %d });
  await page.goto('%s');

  for (const p of plugins) {
    if (typeof p.onPageLoad === 'function') await p.onPageLoad(page);
  }

  const title = await page.title();
  console.log("[Phantom Vite] Title:", title);
  await page.screenshot({ path: 'screenshot.png' });
  await browser.close();

  for (const p of plugins) {
    if (typeof p.onExit === 'function') await p.onExit();
  }
})();
`, cfg.Headless, cfg.Viewport.Width, cfg.Viewport.Height, url)

	return writeTempScriptFile(script, "puppet-temp.mjs")
}

func writePlaywrightScript(url string) (string, error) {
	cfg := loadConfig()

	script := fmt.Sprintf(`import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch({ headless: %v });
  const context = await browser.newContext({ viewport: { width: %d, height: %d } });
  const page = await context.newPage();
  await page.goto('%s');

  const title = await page.title();
  console.log("[Phantom Vite] Title:", title);
  await page.screenshot({ path: 'screenshot.png' });
  await browser.close();
})();
`, cfg.Headless, cfg.Viewport.Width, cfg.Viewport.Height, url)

	return writeTempScriptFile(script, "playwright-temp.mjs")
}

func writeSeleniumScript(url string) (string, error) {
	cfg := loadConfig()

	script := fmt.Sprintf(`# selenium_temp.py
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
import time

options = Options()
options.headless = %v
options.add_argument("--window-size=%d,%d")
driver = webdriver.Chrome(options=options)

driver.get("%s")
print("[Phantom Vite] Title:", driver.title)
driver.save_screenshot("screenshot.png")
driver.quit()
`, cfg.Headless, cfg.Viewport.Width, cfg.Viewport.Height, url)

	return writeTempScriptFile(script, "selenium-temp.py")
}

func writeTempScriptFile(content, filename string) (string, error) {
	tempPath := filepath.Join(os.TempDir(), filename)
	err := os.WriteFile(tempPath, []byte(content), 0644)
	if err != nil {
		return "", err
	}
	return tempPath, nil
}
