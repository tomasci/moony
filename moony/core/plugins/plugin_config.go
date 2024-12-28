package plugins

type PluginConfig struct {
	Name         string      `json:"name"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	Version      string      `json:"version"`
	Author       string      `json:"author"`
	Dependencies interface{} `json:"dependencies"`
}
