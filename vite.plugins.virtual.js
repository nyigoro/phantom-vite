import fs from 'fs'

export function virtualPluginLoader(configPath = './phantomvite.config.json') {
  const virtualId = 'virtual:phantom-plugins'
  const resolvedVirtualId = '\0' + virtualId

  return {
    name: 'phantom-vite-virtual-plugin-loader',

    resolveId(id) {
      if (id === virtualId) {
        return resolvedVirtualId
      }
    },

    load(id) {
      if (id === resolvedVirtualId) {
        try {
          const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'))
          const plugins = Array.isArray(config.plugins) ? config.plugins : []

          const imports = plugins
            .map((p, i) => `import * as plugin${i} from '${p}'`)
            .join('\n')

          const exports = `export default [${plugins.map((_, i) => `plugin${i}`).join(', ')}]`

          return `${imports}\n\n${exports}`
        } catch (e) {
          console.warn('[Phantom Vite] Failed to load plugins for virtual module:', e)
          return `export default []`
        }
      }
    }
  }
}
