// vite.config.js
import { defineConfig } from 'vite'
import fs from 'fs'
import path from 'path'

const configPath = './phantomvite.config.json'
let entries = ['scripts/example.ts'] // default fallback

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
  console.warn('[Phantom Vite] Failed to parse config:', e)
}

const input = {}
entries.forEach(entry => {
  const name = path.basename(entry, path.extname(entry))
  if (fs.existsSync(entry)) {
    input[name] = path.resolve(entry)
  } else {
    console.warn(`[Phantom Vite] ⚠️ Missing entry: ${entry}`)
  }
})

// Built-in Node.js modules (exclude from bundling)
const NODE_BUILTINS = [
  'fs', 'path', 'os', 'url', 'child_process', 'stream', 'crypto',
  'net', 'tls', 'readline', 'util', 'events', 'buffer',
  'node:fs', 'node:path', 'node:os', 'node:url', 'node:child_process',
  'node:stream', 'node:crypto', 'node:net', 'node:tls', 'node:readline',
  'node:util', 'node:events', 'node:buffer'
]

export default defineConfig({
  build: {
    target: 'node22',
    lib: {
      entry: input,
      formats: ['es'],
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
