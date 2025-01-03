package events

import (
	"context"
	"errors"
	"log"
	"moony/moony/core/dispatcher"
	"moony/moony/core/plugins"
	"moony/moony/utils/response"
	"net"
)

type EventProps struct {
	eventCtx context.Context
	conn     *net.UDPConn
	address  *net.UDPAddr
}
type EventCreateHandler func(data []any, eventProps EventProps)

func Create(pluginConfig plugins.PluginConfig, eventName string, eventHandler EventCreateHandler) {
	d := dispatcher.GetGlobalDispatcher()
	pluginEventName := pluginConfig.Name + "_" + eventName

	d.RegisterEventHandler(pluginEventName, func(eventCtx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any) {
		eventProps := EventProps{
			eventCtx: eventCtx,
			conn:     conn,
			address:  address,
		}

		eventHandler(data, eventProps)
	})

	log.Println("Registered event:", pluginEventName, pluginEventName+"_result", pluginEventName+"_error")
}

func SendError(pluginConfig plugins.PluginConfig, eventName string, errorMessageCode string, eventProps EventProps) {
	d := dispatcher.GetGlobalDispatcher()
	pluginEventName := pluginConfig.Name + "_" + eventName + "_error"

	log.Println(pluginEventName, errorMessageCode)

	if eventProps.conn != nil && eventProps.address != nil {
		response.SendResponse[any](eventProps.conn, eventProps.address, pluginConfig.Name, eventName+"_error", nil, errors.New(errorMessageCode))
	}

	d.Dispatch(pluginEventName, eventProps.eventCtx, eventProps.conn, eventProps.address, []any{errorMessageCode})
}

func Send(pluginConfig plugins.PluginConfig, eventName string, data []any, eventProps EventProps) {
	d := dispatcher.GetGlobalDispatcher()
	pluginEventName := pluginConfig.Name + "_" + eventName + "_result"

	if eventProps.conn != nil && eventProps.address != nil {
		response.SendResponse(eventProps.conn, eventProps.address, pluginConfig.Name, eventName+"_result", data, nil)
	}

	d.Dispatch(pluginEventName, eventProps.eventCtx, eventProps.conn, eventProps.address, data)
}

func BroadcastError() {
	// todo: add broadcasting
	return
}

func Broadcast() {
	// todo: add broadcasting
	return
}
