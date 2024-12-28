package main

import (
	"context"
	"moony/moony/core/event_dispatcher"
	"moony/moony/core/plugins"
	"moony/moony/utils/response"
	"net"
)

type PingPlugin struct {
	config plugins.PluginConfig
}

func init() {
	return
}

func (plugin *PingPlugin) Init(ctx context.Context, dispatcher *event_dispatcher.EventDispatcher, config plugins.PluginConfig) error {
	plugin.config = config

	// register plugin command
	dispatcher.RegisterEventHandler(plugin.config.Name+".ping", func(eventCtx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any) {
		responseData := "pong"
		// send response to client
		if conn != nil && address != nil {
			response.SendResponse[any](conn, address, plugin.config.Name, "ping", responseData, nil)
		}
		// notify local listeners
		dispatcher.Dispatch(plugin.config.Name+".ping.result", eventCtx, conn, address, []any{responseData})
	})

	return nil
}

var PluginInstance PingPlugin
