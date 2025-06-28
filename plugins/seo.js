export const onStart = (ctx) => {
  console.log("[SEO Plugin] onStart with context:", ctx);
  if (ctx.engine === 'puppeteer') {
    console.log("Launching Puppeteer!");
  }
};

export async function onPageLoad(page) {
  console.log("[SEO Plugin] onPageLoad: Checking SEO elements...");

  const title = await page.title();
  const description = await page.$eval('meta[name="description"]', el => el.content).catch(() => null);
  const h1 = await page.$eval('h1', el => el.innerText).catch(() => null);

  console.log(`[SEO Plugin] Title: ${title}`);
  console.log(`[SEO Plugin] Meta description: ${description || "Not found"}`);
  console.log(`[SEO Plugin] First <h1>: ${h1 || "Not found"}`);
}

export async function onExit() {
  console.log("[SEO Plugin] onExit: Analysis complete.");
}
