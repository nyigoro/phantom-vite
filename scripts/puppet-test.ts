import puppeteer from 'puppeteer';

(async () => {
  console.log("[Phantom Vite] Launching headless browser...");

  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();

  const url = 'https://example.com';
  await page.goto(url);

  const title = await page.title();
  console.log(`[Phantom Vite] Page title: "${title}"`);

  // Optional screenshot
  await page.screenshot({ path: 'example.png' });
  console.log("[Phantom Vite] Screenshot saved: example.png");

  await browser.close();
})();
