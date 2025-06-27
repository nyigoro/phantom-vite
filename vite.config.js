// vite.config.js
import { defineConfig } from 'vite'
import fs from 'fs'
import path from 'path'

const configPath = './phantomvite.config.json'
let entries = ['scripts/example.ts'] // fallback

try {
  if (fs.existsSync(configPath)) {
    const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'))
    if (Array.isArray(config.entries)) entries = config.entries
    else if (Array.isArray(config.entry)) entries = config.entry
    else if (typeof config.entry === 'string') entries = [config.entry]
  }
} catch (e) {
  console.warn('[Phantom Vite] Failed to parse config:', e)
}

// Validate existence
const validEntries = []
entries.forEach(entry => {
  if (fs.existsSync(entry)) validEntries.push(entry)
  else console.warn(`❌ Missing: ${entry}`)
})
if (validEntries.length === 0) {
  console.error('[Phantom Vite] No valid entries found.')
  process.exit(1)
}

// Map entries to { [name]: path }
const input = {}
validEntries.forEach(entry => {
  const name = path.basename(entry, path.extname(entry))
  input[name] = entry
})

// Toggle bundling
const shouldBundle = process.env.BUNDLE === 'true'

// Node.js built-ins
const NODE_BUILTINS = [
  'fs', 'path', 'http', 'https', 'url', 'util', 'os', 'dns',
  'net', 'tls', 'stream', 'buffer', 'crypto', 'events',
  'child_process', 'readline', 'worker_threads',
  // Node with 'node:' prefix
  ...[
    'fs', 'path', 'http', 'https', 'url', 'util', 'os', 'dns',
    'net', 'tls', 'stream', 'buffer', 'crypto', 'events',
    'child_process', 'readline', 'worker_threads'
  ].map(m => `node:${m}`)
]

export default defineConfig({
  build: {
    target: 'node22',
    lib: {
      entry: input,
      formats: ['es']
    },
    outDir: 'dist',
    emptyOutDir: false,
    minify: false,
    sourcemap: true,
    rollupOptions: {
      external: shouldBundle ? [] : (id => {
        return NODE_BUILTINS.includes(id) ||
               /^puppeteer/.test(id) ||
               /^playwright/.test(id) ||
               /^selenium/.test(id)
      }),
      output: {
        format: 'es',
        entryFileNames: '[name].js',
        preserveModules: false
      }
    }
  },
  define: {
    global: 'globalThis'
  },
  optimizeDeps: {
    exclude: ['puppeteer', 'puppeteer-core', '@puppeteer/browsers']
  }
})
