import puppeteer from 'puppeteer'

;(async () => {
  const browser = await puppeteer.launch({ headless: true })
  const page = await browser.newPage()
  await page.goto('https://example.com')
  const title = await page.title()
  console.log('[Puppeteer] Page title:', title)
  await browser.close()
})()
