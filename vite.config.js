// vite.config.js
import { defineConfig } from 'vite'
import fs from 'fs'
import path from 'path'

// Load Phantom Vite config
const configPath = './phantomvite.config.json'
let entries = ['scripts/example.ts'] // default

if (fs.existsSync(configPath)) {
  const cfg = JSON.parse(fs.readFileSync(configPath, 'utf-8'))
  if (Array.isArray(cfg.entries)) entries = cfg.entries
  else if (Array.isArray(cfg.entry)) entries = cfg.entry
  else if (typeof cfg.entry === 'string') entries = [cfg.entry]
}

// Validate & filter only .ts/.js
entries = entries.filter(file =>
  fs.existsSync(file) && /\.(ts|js)$/.test(file)
)

if (entries.length === 0) {
  console.error('[Phantom Vite] No valid .ts or .js entry files found!')
  process.exit(1)
}

const input = {}
entries.forEach(entry => {
  const name = path.basename(entry, path.extname(entry))
  input[name] = path.resolve(entry)
})

export default defineConfig({
  build: {
    ssr: true, // ✅ Target Node.js explicitly
    target: 'node20',
    outDir: 'dist',
    emptyOutDir: false,
    sourcemap: true,
    minify: false,

    lib: {
      entry: input,
      formats: ['es'],
    },

    rollupOptions: {
      input,
      output: {
        format: 'es',
        entryFileNames: '[name].js',
      },
      external: [
        // Node built-ins
        'fs', 'fs/promises', 'path', 'url', 'os', 'crypto', 'stream',
        'child_process', 'events', 'buffer', 'util', 'http', 'https',
        'readline', 'zlib', 'assert', 'dns', 'net', 'tls', 'module',
        'node:fs', 'node:fs/promises', 'node:path', 'node:url',
        'node:os', 'node:crypto', 'node:stream', 'node:child_process',
        'node:events', 'node:buffer', 'node:util', 'node:http', 'node:https',
        'node:readline', 'node:zlib', 'node:assert', 'node:dns',
        'node:net', 'node:tls', 'node:module',

        // Major automation deps
        'puppeteer', 'puppeteer-core', '@puppeteer/browsers', 'playwright',
        'proxy-agent', 'get-uri', 'pac-resolver', 'basic-ftp',
        '@tootallnate/quickjs-emscripten', 'smart-buffer'
      ],
    },
  },

  // Prevent Vite from pre-optimizing Node modules
  optimizeDeps: {
    exclude: [
      'puppeteer', 'puppeteer-core', '@puppeteer/browsers', 'playwright'
    ]
  },

  assetsInclude: [
    '**/*.py',
    '**/*.gemini'
  ],
  define: {
    global: 'globalThis',
  },
})
