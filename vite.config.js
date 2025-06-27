// vite.config.js - Multi-entry version
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
    } else if (typeof config.entry === 'string') {
      entries = [config.entry]
    } else if (Array.isArray(config.entry)) {
      entries = config.entry
    }
  }
} catch (e) {
  console.warn('[Phantom Vite] Failed to parse config:', e)
}

// Build input object for multiple entries
const input = {}
entries.forEach(entry => {
  const name = path.basename(entry, path.extname(entry))
  input[name] = entry
})

console.log('[Phantom Vite] Building entries:', Object.keys(input))

export default defineConfig({
  build: {
    lib: {
      entry: input,
      formats: ['es']
    },
    outDir: 'dist',
    emptyOutDir: false,
    minify: false,
    rollupOptions: {
      output: {
        entryFileNames: '[name].js'
      }
    }
  }
})
