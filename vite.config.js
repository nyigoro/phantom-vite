import { defineConfig } from 'vite'
import fs from 'fs'
import path from 'path'
import { virtualPluginLoader } from './vite.plugins.virtual.js' 

const configPath = './phantomvite.config.json'
let entries = ['scripts/example.ts'] 

try {
  if (fs.existsSync(configPath)) {
    const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'))

    // entries > entry[] > entry string
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
  if (fs.existsSync(entry)) {
    const name = path.basename(entry, path.extname(entry))
    input[name] = entry
  } else {
    console.warn(`[Phantom Vite] Missing entry file: ${entry}`)
  }
})

// External Node.js built-ins
const NODE_BUILTINS = [
  'fs', 'path', 'http', 'https', 'url', 'util', 'os', 'dns', 'net', 'tls', 'stream',
  'buffer', 'crypto', 'events', 'child_process', 'readline', 'worker_threads',
  'module',
  'node:fs', 'node:path', 'node:child_process', 'node:fs/promises',
  'node:url', 'node:stream', 'node:readline', 'node:process',
  'node:crypto', 'node:dns', 'node:events', 'node:buffer', 'node:assert'
]

export default defineConfig({
  plugins: [virtualPluginLoader()],
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
    exclude: ['puppeteer', 'puppeteer-core', '@puppeteer/browsers']
  }
})
