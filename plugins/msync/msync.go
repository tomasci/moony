package main

import (
	"context"
	"moony/moony/core/plugins"
)

type MSyncPlugin struct {
	config plugins.PluginConfig
}

func (plugin *MSyncPlugin) Init(ctx context.Context, config plugins.PluginConfig) error {
	return nil
}

var PluginInstance MSyncPlugin
