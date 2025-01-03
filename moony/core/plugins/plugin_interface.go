package plugins

import (
	"context"
)

// Plugin - an interface all plugins must implement to work with server
type Plugin interface {
	Init(ctx context.Context, config PluginConfig) error
}
