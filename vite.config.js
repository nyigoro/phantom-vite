// vite.config.js
import { defineConfig } from 'vite'
import fs from 'fs'
import path from 'path'

const configPath = './phantomvite.config.json'
let entries = ['scripts/example.ts'] // fallback if none provided

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

// Filter valid files
const input = {}
entries.forEach(entry => {
  if (fs.existsSync(entry)) {
    const name = path.basename(entry, path.extname(entry))
    input[name] = path.resolve(entry)
  } else {
    console.warn(`[Phantom Vite] Entry not found: ${entry}`)
  }
})

// ❗ Node.js built-ins (do NOT bundle these)
const NODE_BUILTINS = [
  'fs', 'path', 'os', 'url', 'http', 'https', 'net', 'tls', 'child_process',
  'stream', 'buffer', 'util', 'events', 'readline',
  'node:fs', 'node:path', 'node:os', 'node:url', 'node:http', 'node:https',
  'node:net', 'node:tls', 'node:stream', 'node:buffer', 'node:util', 'node:events'
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
      // ✅ Only exclude Node built-ins — all other npm deps are bundled
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
  }
})
