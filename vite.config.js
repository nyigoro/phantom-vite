// vite.config.js
import { defineConfig } from 'vite'
import fs from 'fs'
import path from 'path'

const configPath = './phantomvite.config.json'
let entries = {}

// Parse entry files
try {
  if (fs.existsSync(configPath)) {
    const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'))

    // Handle single string or array of entries
    const rawEntries = typeof config.entry === 'string' ? [config.entry] : config.entry
    if (Array.isArray(rawEntries)) {
      for (const entryPath of rawEntries) {
        const baseName = path.basename(entryPath, path.extname(entryPath))
        entries[baseName] = path.resolve(__dirname, entryPath)
      }
    }
  }
} catch (e) {
  console.warn('[Phantom Vite] Failed to parse phantomvite.config.json:', e)
}

if (Object.keys(entries).length === 0) {
  // Default fallback
  entries = {
    example: path.resolve(__dirname, 'scripts/example.ts'),
  }
}

export default defineConfig({
  build: {
    target: 'node20',
    outDir: 'dist',
    emptyOutDir: false,
    minify: false,
    sourcemap: true,

    rollupOptions: {
      input: entries,
      output: {
        format: 'es',
        entryFileNames: '[name].js',
        preserveModules: false,
      },
      external: [
        'fs', 'path', 'http', 'https', 'url', 'util', 'os', 'dns', 
        'net', 'tls', 'stream', 'buffer', 'crypto', 'events',
        'child_process', 'readline', 'worker_threads',

        'node:fs', 'node:path', 'node:http', 'node:https', 'node:url', 
        'node:util', 'node:os', 'node:dns', 'node:net', 'node:tls',
        'node:stream', 'node:buffer', 'node:crypto', 'node:events',
        'node:child_process', 'node:readline', 'node:worker_threads',

        'puppeteer', 'puppeteer-core', '@puppeteer/browsers',
        'proxy-agent', 'get-uri', 'pac-resolver', 'basic-ftp',
        '@tootallnate/quickjs-emscripten', 'smart-buffer'
      ],
    },
  },

  define: {
    global: 'globalThis',
  },

  optimizeDeps: {
    exclude: ['puppeteer', 'puppeteer-core', '@puppeteer/browsers']
  }
})
