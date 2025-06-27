// vite.plugins.virtual.js
import fs from 'fs'
import path from 'path'

export function virtualPluginLoader() {
  const configPath = './phantomvite.config.json'
  let pluginPaths = []

  if (fs.existsSync(configPath)) {
    const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'))
    pluginPaths = config.plugins || []
  }

  return {
    name: 'phantomvite-virtual-plugins',
    resolveId(id) {
      if (id === 'virtual:phantom-plugins') return id
    },
    load(id) {
      if (id === 'virtual:phantom-plugins') {
        // Generate imports + plugin array
        const imports = pluginPaths.map((p, i) => {
          const fullPath = path.resolve(p)
          return `import * as plugin${i} from ${JSON.stringify(fullPath)}`
        }).join('\n')

        const pluginList = pluginPaths.map((_, i) => `plugin${i}`).join(', ')

        return `
${imports}

export const plugins = [${pluginList}]
`
      }
    }
  }
}
