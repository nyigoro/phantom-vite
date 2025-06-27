// vite.config.js - Node.js compatible build
import { defineConfig } from 'vite'
import fs from 'fs'

let entryFile = 'scripts/example.ts'
const configPath = './phantomvite.config.json'

try {
  if (fs.existsSync(configPath)) {
    const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'))
    
    // Handle both string and array entries
    if (typeof config.entry === 'string') {
      entryFile = config.entry
    } else if (Array.isArray(config.entry) && config.entry.length > 0) {
      entryFile = config.entry[0]
      console.log(`[Phantom Vite] Using first entry: ${entryFile}`)
    }
  }
} catch (e) {
  console.warn('[Phantom Vite] Failed to parse config:', e)
}

// Check if entry file exists
if (!fs.existsSync(entryFile)) {
  console.error(`[Phantom Vite] Entry file not found: ${entryFile}`)
  console.log('[Phantom Vite] Available files in scripts/:')
  try {
    const scriptsDir = 'scripts'
    if (fs.existsSync(scriptsDir)) {
      const files = fs.readdirSync(scriptsDir)
      files.forEach(file => console.log(`  - ${scriptsDir}/${file}`))
    }
  } catch (e) {
    console.log('  (could not read scripts directory)')
  }
}

export default defineConfig({
  build: {
    // Target Node.js environment
    target: 'node20',
    
    lib: {
      entry: entryFile,
      name: 'PhantomScript',
      fileName: () => entryFile.split('/').pop().replace(/\.(ts|js)$/, ''),
      formats: ['es']
    },
    
    outDir: 'dist',
    emptyOutDir: false,
    minify: false,
    
    rollupOptions: {
      // Mark Node.js built-ins and large dependencies as external
      external: [
        // Node.js built-ins
        'fs', 'path', 'http', 'https', 'url', 'util', 'os', 'dns', 
        'net', 'tls', 'stream', 'buffer', 'crypto', 'events',
        'child_process', 'readline', 'worker_threads',
        
        // Node.js built-ins with node: prefix
        'node:fs', 'node:path', 'node:http', 'node:https', 'node:url', 
        'node:util', 'node:os', 'node:dns', 'node:net', 'node:tls',
        'node:stream', 'node:buffer', 'node:crypto', 'node:events',
        'node:child_process', 'node:readline', 'node:worker_threads',
        
        // Puppeteer and related packages
        'puppeteer', 'puppeteer-core', '@puppeteer/browsers',
        
        // Other Node.js-only packages
        'proxy-agent', 'get-uri', 'pac-resolver', 'basic-ftp',
        '@tootallnate/quickjs-emscripten', 'smart-buffer'
      ],
      
      output: {
        format: 'es',
        entryFileNames: '[name].js',
        // Preserve module structure for Node.js
        preserveModules: false,
        // Generate sourcemaps for debugging
        sourcemap: true
      }
    }
  },
  
  // Disable browser-specific optimizations
  define: {
    global: 'globalThis',
  },
  
  // Enable Node.js polyfills if needed
  optimizeDeps: {
    // Skip pre-bundling for Node.js dependencies
    exclude: ['puppeteer', 'puppeteer-core', '@puppeteer/browsers']
  }
})
