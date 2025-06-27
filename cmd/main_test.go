package main

import "testing"

func writeTempScript(url string) (string, error) {
	cfg := loadConfig()
	pluginsJSON, _ := json.Marshal(cfg.Plugins)

	code := fmt.Sprintf(`import puppeteer from 'puppeteer';

const pluginPaths = process.env["PHANTOM_PLUGINS"]?.split(",") ?? [];
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

  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();

  const url = '%s';
  await page.goto(url);

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
`, url)
