// vite.config.js - Multi-entry Node.js compatible build
import { defineConfig } from 'vite'
import fs from 'fs'
import path from 'path'

const configPath = './phantomvite.config.json'
let entries = ['scripts/example.ts'] // default fallback

try {
  if (fs.existsSync(configPath)) {
    const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'))
    
    // Priority order: entries array > entry array > entry string
    if (Array.isArray(config.entries)) {
      entries = config.entries
      console.log(`[Phantom Vite] Found ${entries.length} entries in config.entries`)
    } else if (Array.isArray(config.entry)) {
      entries = config.entry
      console.log(`[Phantom Vite] Found ${entries.length} entries in config.entry`)
    } else if (typeof config.entry === 'string') {
      entries = [config.entry]
      console.log(`[Phantom Vite] Found single entry: ${config.entry}`)
    }
  }
} catch (e) {
  console.warn('[Phantom Vite] Failed to parse config:', e)
}

// Validate that all entry files exist
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
  missingEntries.forEach(entry => console.warn(`  ❌ ${entry}`))
}

if (validEntries.length === 0) {
  console.error('[Phantom Vite] No valid entry files found! Check your phantomvite.config.json')
  process.exit(1)
}

console.log('[Phantom Vite] Building entries:')
validEntries.forEach(entry => console.log(`  ✅ ${entry}`))

// Build input object for multiple entries
const input = {}
validEntries.forEach(entry => {
  const name = path.basename(entry, path.extname(entry))
  input[name] = entry
})

// Define Node.js built-in modules that should always be external
const NODE_BUILTINS = [
  'fs', 'path', 'http', 'https', 'url', 'util', 'os', 'dns', 
  'net', 'tls', 'stream', 'buffer', 'crypto', 'events',
  'child_process', 'readline', 'worker_threads',
  // Node.js built-ins with node: prefix
  'node:fs', 'node:path', 'node:http', 'node:https', 'node:url', 
  'node:util', 'node:os', 'node:dns', 'node:net', 'node:tls',
  'node:stream', 'node:buffer', 'node:crypto', 'node:events',
  'node:child_process', 'node:readline', 'node:worker_threads',
];

export default defineConfig({
  build: {
    // Target Node.js environment
    target: 'node22',
    
    lib: {
      entry: input,
      formats: ['es']
    },
    
    outDir: 'dist',
    emptyOutDir: false,
    minify: false, // Keep false for debugging, true for production
    sourcemap: true,
    
    rollupOptions: {
      // ONLY mark Node.js built-ins as external.
      // All other npm dependencies will now be bundled into the output.
      external: NODE_BUILTINS,
      
      output: {
        format: 'es',
        entryFileNames: '[name].js',
        // Preserve module structure for Node.js - this generally means
        // that each entry file will result in its own output file,
        // but its dependencies will be bundled into it or shared chunks.
        preserveModules: false 
      }
    }
  },
  
  // Disable browser-specific optimizations
  define: {
    global: 'globalThis',
  },
  
  // OptimizeDeps is for dev server, not production build, but keep it
  // if you're not seeing issues in dev.
  // We don't need to exclude puppeteer etc here for build, as Rollup's external handles it.
  optimizeDeps: {
    // Skip pre-bundling for Node.js dependencies
    exclude: ['puppeteer', 'puppeteer-core', '@puppeteer/browsers']
  }
})
