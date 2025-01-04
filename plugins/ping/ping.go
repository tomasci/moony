package main

import (
	"context"
	"moony/moony/core/events"
	"moony/moony/core/plugins"
)

type PingPlugin struct {
	config plugins.PluginConfig
}

func init() {
	return
}

func (plugin *PingPlugin) Init(ctx context.Context, config plugins.PluginConfig) error {
	plugin.config = config

	// register plugin command
	events.Create(plugin.config, "ping", func(data []any, eventProps events.EventProps) {
		result := "pong"
		events.Send(plugin.config, "ping", []any{result}, eventProps)
	})

	return nil
}

var PluginInstance PingPlugin
