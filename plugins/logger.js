const allowedCommands = ['open', 'agent'];

export function onStart(context) {
  if (!context) return console.warn('[logger Plugin] No context');

  const { engine, meta } = context;
  const command = meta?.command || 'unknown';
  const url = meta?.url || 'N/A';

  console.log(`[logger Plugin] onStart triggered by command: ${command}`);
  console.log(`[logger Plugin] Engine: ${engine}`);
  console.log(`[logger Plugin] Target URL: ${url}`);
}


