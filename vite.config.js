import { defineConfig } from 'vite'

export default defineConfig({
  build: {
    lib: {
      entry: 'scripts/example.ts',
      name: 'ExampleScript',
      fileName: 'example',
      formats: ['es']
    },
    outDir: 'dist',
    emptyOutDir: false,
    minify: false
  }
})

