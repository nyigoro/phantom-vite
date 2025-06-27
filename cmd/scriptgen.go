func writeTempScript(url, engine string) (string, error) {
    var content string

    switch engine {
    case "puppeteer":
        content = fmt.Sprintf(`import puppeteer from 'puppeteer';
(async () => {
  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();
  await page.goto('%s');
  console.log('[Puppeteer] Title:', await page.title());
  await browser.close();
})();`, url)

    case "playwright":
        content = fmt.Sprintf(`import { chromium } from 'playwright';
(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  await page.goto('%s');
  console.log('[Playwright] Title:', await page.title());
  await browser.close();
})();`, url)

    case "selenium":
        content = fmt.Sprintf(`from selenium import webdriver
from selenium.webdriver.chrome.options import Options

options = Options()
options.add_argument("--headless")
driver = webdriver.Chrome(options=options)
driver.get("%s")
print("[Selenium] Title:", driver.title)
driver.quit()`, url)

    case "gemini":
        content = fmt.Sprintf(`gemini open %s
gemini expect-title "Example Domain"`, url)

    default:
        return "", fmt.Errorf("unsupported engine: %s", engine)
    }

    ext := map[string]string{
        "puppeteer": ".ts",
        "playwright": ".ts",
        "selenium": ".py",
        "gemini": ".gemini",
    }[engine]

    tmpFile, err := os.CreateTemp("", "phantom-"+engine+"-*"+ext)
    if err != nil {
        return "", err
    }
    defer tmpFile.Close()

    _, err = tmpFile.WriteString(content)
    return tmpFile.Name(), err
}
