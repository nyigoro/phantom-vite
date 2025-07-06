const allowedCommands = ['open', 'agent'];

export function onStart(context) {
  if (!context) return console.warn('[SEO Plugin] No context');

  const { engine, meta } = context;
  const command = meta?.command || 'unknown';
  const url = meta?.url || 'N/A';

  console.log(`[SEO Plugin] onStart triggered by command: ${command}`);
  console.log(`[SEO Plugin] Engine: ${engine}`);
  console.log(`[SEO Plugin] Target URL: ${url}`);
}


