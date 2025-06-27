// vite.config.js - Multi-entry Node.js compatible build
import { defineConfig } from 'vite'
import fs from 'fs'
import path from 'path'

const configPath = './phantomvite.config.json'
let entries = ['scripts/example.ts'] 

try {
  if (fs.existsSync(configPath)) {
    const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'))

    if (Array.isArray(config.entries)) {
      entries = config.entries
    } else if (Array.isArray(config.entry)) {
      entries = config.entry
    } else if (typeof config.entry === 'string') {
      entries = [config.entry]
    }
  }
} catch (e) {
  console.warn('[Phantom Vite] Failed to parse phantomvite.config.json:', e)
}

// Validate entry existence
const validEntries = []
const missingEntries = []

entries.forEach(entry => {
  if (fs.existsSync(entry)) {
    validEntries.push(entry)
  } else {
    missingEntries.push(entry)
  }
})

if (missingEntries.length > 0) {
  console.warn('[Phantom Vite] Missing entry files:')
  missingEntries.forEach(e => console.warn(`  ❌ ${e}`))
}

if (validEntries.length === 0) {
  console.error('[Phantom Vite] No valid entry files found!')
  process.exit(1)
}

// Map input
const input = {}
validEntries.forEach(entry => {
  const name = path.basename(entry, path.extname(entry))
  input[name] = path.resolve(__dirname, entry)
})

export default defineConfig({
  assetsInclude: [
    '**/*.py',
    '**/*.gemini'
  ],
  build: {
    target: 'node18',
    outDir: 'dist',
    emptyOutDir: false,
    sourcemap: true,
    minify: false,

    rollupOptions: {
      input,
      output: {
        format: 'es',
        entryFileNames: '[name].js',
        preserveModules: false,
      },
      // 🟡 Leave Puppeteer bundled (do NOT mark as external)
      external: [
        'fs', 'path', 'os', 'util', 'url', 'crypto', 'stream', 'child_process',
        'events', 'buffer', 'readline', 'worker_threads', 'http', 'https',
        'dns', 'net', 'tls',

        'node:fs', 'node:path', 'node:os', 'node:util', 'node:url',
        'node:crypto', 'node:stream', 'node:child_process', 'node:events',
        'node:buffer', 'node:readline', 'node:worker_threads', 'node:http',
        'node:https', 'node:dns', 'node:net', 'node:tls'
        
        // 🔥 We intentionally do NOT externalize puppeteer
      ]
    }
  },

  define: {
    global: 'globalThis'
  },

  optimizeDeps: {
    // Prevent pre-bundling Puppeteer (just in case)
    exclude: ['puppeteer']
  }
})
