package plugins

import "sync"

var (
	// for thread-safe read-write access to plugins
	pluginsMu sync.RWMutex
	// plugin list
	plugins = make(map[string]Plugin)
)

// RegisterPlugin - use to add new plugin to list
func RegisterPlugin(name string, plugin Plugin) {
	// lock plugin list
	pluginsMu.Lock()
	// unlock when done
	defer pluginsMu.Unlock()
	// add new plugin under provided name
	plugins[name] = plugin
}

// GetPlugin - get plugin by plugin name
func GetPlugin(name string) (Plugin, bool) {
	// lock plugin list, so no new plugins can be added
	pluginsMu.RLock()
	// unlock when done
	defer pluginsMu.RUnlock()
	// get plugin & check if it exists
	plugin, exists := plugins[name]
	// return plugin and status
	return plugin, exists
}
