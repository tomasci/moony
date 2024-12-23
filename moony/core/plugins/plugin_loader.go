package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"moony/moony/core/event_dispatcher"
	"os"
	"path/filepath"
	"plugin"
)

// LoadPlugins - use it to load all plugins from specified directory
func LoadPlugins(pluginDir string, dispatcher *event_dispatcher.EventDispatcher) (int, error) {
	// read provided plugins directory
	dirs, err := os.ReadDir(pluginDir)
	if err != nil {
		return 0, fmt.Errorf("error reading plugin directory %s: %s", pluginDir, err)
	}

	// remember all loaded plugins, just in case, maybe it will be used later
	var loadedPlugins []Plugin

	// go through files and folders in directory
	for _, dir := range dirs {
		// skip files, because plugins must be placed in their own directories
		// it must be this way to avoid Plugin Directory Hell (it is when each plugin has its own files, configs, etc.)
		if !dir.IsDir() {
			continue
		}

		// path to plugin directory
		pluginPath := filepath.Join(pluginDir, dir.Name())
		// path to plugins Shared Object (.so) file
		soFilePath := filepath.Join(pluginPath, dir.Name()+".so")
		// path to plugin configuration
		pluginJsonPath := filepath.Join(pluginPath, "plugin.json")

		// now check everything
		// check if .so file exists
		if _, err := os.Stat(soFilePath); os.IsNotExist(err) {
			log.Println("plugin not found:", soFilePath)
			continue
		}

		// check if plugin.json config exists
		if _, err := os.Stat(pluginJsonPath); os.IsNotExist(err) {
			log.Println("plugin config not found:", pluginJsonPath)
			continue
		}

		// try to load so file
		// yes, you already know if it exists or not, but it was only for logging purposes, and you cannot return fatal here
		// because first - you want to load all other plugins and second - you want to print all problems to users, not one by one
		p, err := plugin.Open(soFilePath)
		if err != nil {
			log.Println("failed to open plugin:", pluginPath, err)
		}

		// look for PluginInstance symbol
		symbol, err := p.Lookup("PluginInstance")
		if err != nil {
			log.Println("failed to find PluginInstance:", pluginPath, err)
			continue
		}

		// check is it is of Plugin type
		plug, ok := symbol.(Plugin)
		if !ok {
			log.Println("symbol PluginInstance is not of type Plugin:", pluginPath)
			continue
		}

		// read plugin.json
		jsonData, err := os.ReadFile(pluginJsonPath)
		if err != nil {
			log.Println("failed to read plugin.json:", pluginPath, err)
		}

		// parse plugin.json
		var config PluginConfig
		if err := json.Unmarshal(jsonData, &config); err != nil {
			log.Println("failed to parse plugin.json:", pluginPath, err)
			continue
		}

		// initialize plugin
		if err := plug.Init(context.Background(), dispatcher, config); err != nil {
			log.Println("failed to init plugin:", pluginPath, err)
			continue
		}

		// register plugin
		RegisterPlugin(config.Name, plug)
		loadedPlugins = append(loadedPlugins, plug)

		log.Printf("Loaded %s@%s by %s\n", config.Name, config.Version, config.Author)
	}

	// return amount of loaded plugins (just for now)
	return len(loadedPlugins), nil
}
