package plugins

import (
	"context"
	"moony/moony/core/event_dispatcher"
)

// Plugin - an interface all plugins must implement to work with server
type Plugin interface {
	Init(ctx context.Context, dispatcher *event_dispatcher.EventDispatcher, config PluginConfig) error
}
