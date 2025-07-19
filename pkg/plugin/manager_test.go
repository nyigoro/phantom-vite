package plugin

import (
	"context"
	"testing"
)

// MockPlugin for testing the PluginManager
type MockPlugin struct {
	name string
	version string
	dependencies []string
}

func (m *MockPlugin) Name() string {
	return m.name
}

func (m *MockPlugin) Version() string {
	return m.version
}

func (m *MockPlugin) Execute(ctx context.Context, engine Engine, args []string) error {
	return nil // Not testing execution here
}

func (m *MockPlugin) Dependencies() []string {
	return m.dependencies
}

// MockEngine for testing the PluginManager
type MockEngine struct {}

func (m *MockEngine) GetConfig() map[string]interface{} {
	return nil // Not testing config retrieval here
}

func (m *MockEngine) Log(message string) {
	// Not testing logging here
}

func TestNewPluginManager(t *testing.T) {
	pm := NewPluginManager()
	if pm == nil {
		t.Errorf("NewPluginManager returned nil")
	}
	if pm.plugins == nil {
		t.Errorf("NewPluginManager did not initialize plugins map")
	}
}

func TestRegisterPlugin(t *testing.T) {
	pm := NewPluginManager()
	plugin := &MockPlugin{name: "testPlugin", version: "1.0"}
	pm.RegisterPlugin(plugin)

	if _, exists := pm.plugins["testPlugin"]; !exists {
		t.Errorf("RegisterPlugin failed to register plugin")
	}
}

func TestGetPlugin(t *testing.T) {
	pm := NewPluginManager()
	plugin := &MockPlugin{name: "testPlugin", version: "1.0"}
	pm.RegisterPlugin(plugin)

	retrievedPlugin, exists := pm.GetPlugin("testPlugin")
	if !exists {
		t.Errorf("GetPlugin failed to retrieve existing plugin")
	}
	if retrievedPlugin.Name() != "testPlugin" {
		t.Errorf("GetPlugin retrieved wrong plugin: expected %s, got %s", "testPlugin", retrievedPlugin.Name())
	}

	_, exists = pm.GetPlugin("nonExistentPlugin")
	if exists {
		t.Errorf("GetPlugin retrieved non-existent plugin")
	}
}

func TestListPlugins(t *testing.T) {
	pm := NewPluginManager()
	plugin1 := &MockPlugin{name: "plugin1", version: "1.0"}
	plugin2 := &MockPlugin{name: "plugin2", version: "1.0"}

	pm.RegisterPlugin(plugin1)
	pm.RegisterPlugin(plugin2)

	plugins := pm.ListPlugins()

	if len(plugins) != 2 {
		t.Errorf("ListPlugins returned incorrect number of plugins: expected %d, got %d", 2, len(plugins))
	}

	found1 := false
	found2 := false
	for _, name := range plugins {
		if name == "plugin1" {
			found1 = true
		} else if name == "plugin2" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("ListPlugins did not return all registered plugins")
	}
}
