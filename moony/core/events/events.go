package events

import (
	"context"
	"errors"
	"log"
	"moony/database/redis"
	"moony/moony/core/dispatcher"
	"moony/moony/core/plugins"
	"moony/moony/utils/response"
	"net"
)

type EventProps struct {
	EventCtx context.Context
	Conn     *net.UDPConn
	Address  *net.UDPAddr
}
type EventCreateHandler func(data []any, eventProps EventProps)

func Create(pluginConfig plugins.PluginConfig, eventName string, eventHandler EventCreateHandler) {
	// get dispatcher
	d := dispatcher.GetGlobalDispatcher()
	// create full event name from plugin and event name, underscore is required
	pluginEventName := pluginConfig.Name + "_" + eventName

	// create event handler
	d.RegisterEventHandler(pluginEventName, func(eventCtx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any) {
		eventProps := EventProps{
			EventCtx: eventCtx,
			Conn:     conn,
			Address:  address,
		}

		eventHandler(data, eventProps)
	})

	// some logs
	log.Println("Registered event:", pluginEventName, pluginEventName+"_result", pluginEventName+"_error")
}

func SendError(pluginConfig plugins.PluginConfig, eventName string, errorMessageCode string, eventProps EventProps) {
	d := dispatcher.GetGlobalDispatcher()
	pluginEventName := pluginConfig.Name + "_" + eventName + "_error"

	log.Println(pluginEventName, errorMessageCode)

	if eventProps.Conn != nil && eventProps.Address != nil {
		// send response to client
		response.SendResponse[any](eventProps.Conn, eventProps.Address, pluginConfig.Name, eventName+"_error", nil, errors.New(errorMessageCode))
	}

	// notify all local listeners
	d.Dispatch(pluginEventName, eventProps.EventCtx, eventProps.Conn, eventProps.Address, []any{errorMessageCode})
}

func Send(pluginConfig plugins.PluginConfig, eventName string, data []any, eventProps EventProps) {
	d := dispatcher.GetGlobalDispatcher()
	pluginEventName := pluginConfig.Name + "_" + eventName + "_result"

	if eventProps.Conn != nil && eventProps.Address != nil {
		response.SendResponse(eventProps.Conn, eventProps.Address, pluginConfig.Name, eventName+"_result", data, nil)
	}

	d.Dispatch(pluginEventName, eventProps.EventCtx, eventProps.Conn, eventProps.Address, data)
}

func BroadcastError() {
	// todo: add broadcasting
	return
}

func Broadcast(pluginConfig plugins.PluginConfig, eventName string, data []any, eventProps EventProps) {
	d := dispatcher.GetGlobalDispatcher()
	pluginEventName := pluginConfig.Name + "_" + eventName + "_result"

	// get redis client
	state, err := redis.GetRedisClient()
	if err != nil {
		log.Println("error getting redis client", err)
		return
	}

	// get connection list with client addresses
	connectionList, err := state.SMembers(eventProps.EventCtx, "connections").Result()
	if err != nil {
		log.Println("error getting connection list", err)
		return
	}

	//log.Println("connection list:", connectionList)

	// send response to everyone in the list
	for _, connection := range connectionList {
		// convert string to actual address
		connectionAddress, err := net.ResolveUDPAddr("udp", connection)
		if err != nil {
			log.Println("error resolving connection address", err)
			return
		}

		if eventProps.Conn != nil {
			response.SendResponse(eventProps.Conn, connectionAddress, pluginConfig.Name, eventName+"_result", data, nil)
		}

		d.Dispatch(pluginEventName, eventProps.EventCtx, eventProps.Conn, connectionAddress, data)
	}
}
