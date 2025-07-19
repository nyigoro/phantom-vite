package plugin

import (
    "context"
)

type PluginManager struct {
    plugins map[string]Plugin
}

type Plugin interface {
    Name() string
    Version() string
    Execute(ctx context.Context, engine Engine, args []string) error
    Dependencies() []string
}

// Engine interface that plugins will interact with
type Engine interface {
    // Define your engine methods here
    GetConfig() map[string]interface{}
    Log(message string)
    // Add other methods as needed
}

// NewPluginManager creates a new plugin manager instance
func NewPluginManager() *PluginManager {
    return &PluginManager{
        plugins: make(map[string]Plugin),
    }
}

// RegisterPlugin adds a plugin to the manager
func (pm *PluginManager) RegisterPlugin(plugin Plugin) {
    pm.plugins[plugin.Name()] = plugin
}

// GetPlugin retrieves a plugin by name
func (pm *PluginManager) GetPlugin(name string) (Plugin, bool) {
    plugin, exists := pm.plugins[name]
    return plugin, exists
}

// ListPlugins returns all registered plugin names
func (pm *PluginManager) ListPlugins() []string {
    names := make([]string, 0, len(pm.plugins))
    for name := range pm.plugins {
        names = append(names, name)
    }
    return names
}