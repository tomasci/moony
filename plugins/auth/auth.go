package main

import (
	"context"
	"moony/moony/core/event_dispatcher"
	"moony/moony/core/plugins"
)

type AuthPlugin struct {
	config plugins.PluginConfig
}

func init() {
	return
}

func (plugin *AuthPlugin) Init(ctx context.Context, dispatcher *event_dispatcher.EventDispatcher, config plugins.PluginConfig) error {
	return nil
}

func (plugin *AuthPlugin) login() {
	return
}

func (plugin *AuthPlugin) create() {
	return
}

var PluginInstance AuthPlugin
