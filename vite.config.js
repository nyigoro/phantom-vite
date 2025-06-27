// vite.config.js
import { defineConfig } from 'vite'
import fs from 'fs'

let entryFile = 'scripts/example.ts'
const configPath = './phantomvite.config.json'

try {
  if (fs.existsSync(configPath)) {
    const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'))
    if (typeof config.entry === 'string') {
      entryFile = config.entry
    }
  }
} catch (e) {
  console.warn('[Phantom Vite] Failed to parse config:', e)
}

export default defineConfig({
  build: {
    lib: {
      entry: entryFile,
      name: 'PhantomScript',
      fileName: () => entryFile.split('/').pop().replace('.ts', ''),
      formats: ['es']
    },
    outDir: 'dist',
    emptyOutDir: false,
    minify: false
  }
})
