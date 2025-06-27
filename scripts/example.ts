// scripts/example.ts
import puppeteer from 'puppeteer'

async function main() {
  console.log('🚀 Starting Phantom Vite example...')
  
  try {
    const browser = await puppeteer.launch({ 
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox']
    })
    
    const page = await browser.newPage()
    await page.goto('https://example.com')
    
    const title = await page.title()
    console.log(`📄 Page title: ${title}`)
    
    // Take a screenshot
    await page.screenshot({ 
      path: 'example-screenshot.png',
      fullPage: true
    })
    console.log('📸 Screenshot saved as example-screenshot.png')
    
    await browser.close()
    console.log('✅ Example completed successfully!')
    
  } catch (error) {
    console.error('❌ Error:', error.message)
    process.exit(1)
  }
}

// Run if this is the main module
if (import.meta.url === `file://${process.argv[1]}`) {
  main()
}

export default main
