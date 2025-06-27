import { defineConfig } from 'vite'
import path from 'path'

const entryFile = process.env.PHANTOM_ENTRY || 'scripts/example.ts'
const name = path.basename(entryFile, path.extname(entryFile))

export default defineConfig({
  build: {
    target: 'node22',
    lib: {
      entry: entryFile,
      name: 'PhantomScript',
      fileName: () => `${name}`,
      formats: ['es']
    },
    outDir: 'dist',
    emptyOutDir: false,
    sourcemap: true,
    minify: false,
    rollupOptions: {
      external: [
        'puppeteer', 'puppeteer-core', 'playwright',
        'fs', 'path', 'os', 'child_process',
        'node:fs', 'node:path', 'node:os', 'node:child_process'
      ],
      output: {
        entryFileNames: '[name].js',
        format: 'es',
        preserveModules: false
      }
    }
  },
  define: {
    global: 'globalThis'
  },
  optimizeDeps: {
    exclude: ['puppeteer', 'playwright']
  }
})
