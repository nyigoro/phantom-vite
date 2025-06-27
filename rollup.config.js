// rollup.config.js
import fs from 'fs'
import path from 'path'
import { defineConfig } from 'rollup'
import commonjs from '@rollup/plugin-commonjs'
import nodeResolve from '@rollup/plugin-node-resolve'
import esbuild from 'rollup-plugin-esbuild'
import json from '@rollup/plugin-json'

const configPath = './phantomvite.config.json'

// 🧠 Load entry files from config
let entries = ['scripts/example.ts'] // fallback
if (fs.existsSync(configPath)) {
  const cfg = JSON.parse(fs.readFileSync(configPath, 'utf-8'))
  if (Array.isArray(cfg.entries)) {
    entries = cfg.entries
  } else if (Array.isArray(cfg.entry)) {
    entries = cfg.entry
  } else if (typeof cfg.entry === 'string') {
    entries = [cfg.entry]
  }
}

// 🔍 Validate entry files
entries = entries.filter(file => {
  if (!fs.existsSync(file)) {
    console.warn(`[Phantom Vite] Missing entry: ${file}`)
    return false
  }
  return true
})

if (entries.length === 0) {
  console.error('[Phantom Vite] No valid entry files found. Exiting.')
  process.exit(1)
}

// 🎯 Build input object
const input = {}
entries.forEach(file => {
  const name = path.basename(file, path.extname(file))
  input[name] = path.resolve(file)
})

// 🚫 Node.js built-ins (won't be bundled)
const nodeBuiltins = [
  'fs', 'path', 'url', 'os', 'crypto', 'stream',
  'child_process', 'events', 'buffer', 'util', 'http', 'https',
  'readline', 'zlib', 'assert', 'dns', 'net', 'tls', 'module',
  'node:fs', 'node:path', 'node:url', 'node:os', 'node:crypto',
  'node:stream', 'node:child_process', 'node:events', 'node:buffer',
  'node:util', 'node:http', 'node:https', 'node:readline',
  'node:zlib', 'node:assert', 'node:dns', 'node:net', 'node:tls', 'node:module'
]

// 🧱 Rollup config
export default defineConfig({
  input,
  output: {
    dir: 'dist',
    format: 'es',
    entryFileNames: '[name].js',
    sourcemap: false,
  },
  external: nodeBuiltins,
  plugins: [
    nodeResolve({
      preferBuiltins: true,
      exportConditions: ['node']
    }),
    commonjs(),
    json(),
    esbuild({
      target: 'node22',
      platform: 'node'
    }),
  ]
})
