export async function onStart() {
  console.log("[Plugin] onStart");
}

export async function onPageLoad(page) {
  console.log("[Plugin] onPageLoad");
  const url = page.url();
  console.log(`[Plugin] Current URL: ${url}`);
}

export async function onExit() {
  console.log("[Plugin] onExit");
}
