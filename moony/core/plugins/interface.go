package plugins

import (
	"context"
)

type PluginConfig struct {
	Name         string      `json:"name"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	Version      string      `json:"version"`
	Author       string      `json:"author"`
	Dependencies interface{} `json:"dependencies"`
}

// Plugin - an interface all plugins must implement to work with server
type Plugin interface {
	Init(ctx context.Context, config PluginConfig) error
}
