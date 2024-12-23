package plugins

type PluginConfig struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Version      string      `json:"version"`
	Author       string      `json:"author"`
	Dependencies interface{} `json:"dependencies"`
}
