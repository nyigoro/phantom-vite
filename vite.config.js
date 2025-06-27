// vite.config.js
import { defineConfig } from 'vite'
import fs from 'fs'
import path from 'path'

// Load Phantom Vite config
const phantomConfigPath = './phantomvite.config.json'
let entries = ['scripts/example.ts'] // default fallback

try {
  if (fs.existsSync(phantomConfigPath)) {
    const config = JSON.parse(fs.readFileSync(phantomConfigPath, 'utf-8'))
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

// Validate entries
const validEntries = []
const missingEntries = []
entries.forEach(entry => {
  if (fs.existsSync(entry)) {
    validEntries.push(entry)
  } else {
    missingEntries.push(entry)
  }
})
if (missingEntries.length) {
  console.warn('[Phantom Vite] Missing entry files:', missingEntries)
}
if (!validEntries.length) {
  console.error('[Phantom Vite] No valid entry files found. Aborting build.')
  process.exit(1)
}

// Rollup input object
const input = {}
validEntries.forEach(entry => {
  const name = path.basename(entry, path.extname(entry))
  input[name] = entry
})

// Node.js built-ins we always want to exclude from bundling
const NODE_BUILTINS = [
  'fs', 'path', 'http', 'https', 'url', 'os', 'util', 'dns',
  'net', 'tls', 'stream', 'zlib', 'assert', 'crypto', 'readline',
  'child_process', 'buffer', 'worker_threads', 'module',
  'node:fs', 'node:path', 'node:http', 'node:https', 'node:url',
  'node:os', 'node:util', 'node:dns', 'node:net', 'node:tls',
  'node:stream', 'node:zlib', 'node:assert', 'node:crypto',
  'node:readline', 'node:child_process', 'node:buffer', 'node:worker_threads',
  'node:module', 'node:process'
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
      external: NODE_BUILTINS,
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
    exclude: ['puppeteer', 'puppeteer-core', '@puppeteer/browsers', 'playwright']
  }
})
