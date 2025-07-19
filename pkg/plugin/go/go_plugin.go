package goplugin

import (
	"fmt"
	"plugin"
)

// Plugin is the interface that Go plugins must implement.
// This is a placeholder and will be expanded as we define the plugin contract.
type Plugin interface {
	onStart(context string)
	onPageLoad(context string)
	onExit(context string)
}

// LoadAndExecuteGoPlugin loads a Go plugin and executes a specified hook.
func LoadAndExecuteGoPlugin(pluginPath, hookName, context string) error {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to load Go plugin %s: %v", pluginPath, err)
	}

	sym, err := p.Lookup(hookName)
	if err != nil {
		// Hook not found, which is acceptable if a plugin doesn't implement all hooks
		return nil
	}

	switch hookName {
	case "onStart":
		if onStartFunc, ok := sym.(func(string)); ok {
			onStartFunc(context)
		} else {
			return fmt.Errorf("onStart hook in %s has wrong signature", pluginPath)
		}
	case "onPageLoad":
		if onPageLoadFunc, ok := sym.(func(string)); ok {
			onPageLoadFunc(context)
		} else {
			return fmt.Errorf("onPageLoad hook in %s has wrong signature", pluginPath)
		}
	case "onExit":
		if onExitFunc, ok := sym.(func(string)); ok {
			onExitFunc(context)
		} else {
			return fmt.Errorf("onExit hook in %s has wrong signature", pluginPath)
		}
	default:
		return fmt.Errorf("unknown hook name: %s", hookName)
	}

	return nil
}
