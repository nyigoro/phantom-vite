// Avoid exporting both in production
// Only keep one active export
export function onStart(context) {
  if (!context) {
    console.warn('[SEO Plugin] No context provided to onStart');
    return;
  }

  if (!context.engine) {
    console.warn('[SEO Plugin] No engine found in context');
    return;
  }

  console.log('[SEO Plugin] Initialized with engine:', context.engine);
}
