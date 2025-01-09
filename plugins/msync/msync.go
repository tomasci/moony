package main

import (
	"context"
	"encoding/json"
	"log"
	"moony/database/redis"
	"moony/moony/core/events"
	"moony/moony/core/plugins"
	"net"
	"time"
)

type MSyncPlugin struct {
	config plugins.PluginConfig
}

type MSyncConnection struct {
	Address   *net.UDPAddr `json:"address"`
	UpdatedAt int64        `json:"updated_at"`
}

const connectionExpirationTime = 5 * time.Second

func (plugin *MSyncPlugin) Init(ctx context.Context, config plugins.PluginConfig) error {
	plugin.config = config

	events.Create(plugin.config, "connect", func(data []any, eventProps events.EventProps) {
		log.Println("Client connected", eventProps.Address)

		err := plugin.onConnect(eventProps)
		if err != nil {
			events.SendError(plugin.config, "connect", "connection_failed", eventProps)
			return
		}

		events.Send(plugin.config, "connect", []any{1}, eventProps)
	})

	events.Create(plugin.config, "disconnect", func(data []any, eventProps events.EventProps) {
		log.Println("Client disconnected", eventProps.Address)

		err := plugin.onDisconnect(eventProps)
		if err != nil {
			// don't need to send error, game is already closed, no one will receive this error
			return
		}

		// same here, no need to send back any response, game already closed
		return
	})

	events.Create(plugin.config, "ping", func(data []any, eventProps events.EventProps) {
		_ = plugin.updateConnection(eventProps)

		result := "pong"
		events.Send(plugin.config, "ping", []any{result}, eventProps)
	})

	return nil
}

func (plugin *MSyncPlugin) onConnect(eventProps events.EventProps) error {
	// get redis client
	state, err := redis.GetRedisClient()
	if err != nil {
		return err
	}

	// client info to store in state
	connection := MSyncConnection{
		Address:   eventProps.Address,
		UpdatedAt: time.Now().UnixMilli(),
	}

	// convert to json
	connectionJsonString, err := json.Marshal(connection)
	if err != nil {
		return err
	}

	// save client info
	err = state.Set(eventProps.EventCtx, "connection:"+eventProps.Address.String(), connectionJsonString, connectionExpirationTime).Err()
	if err != nil {
		return err
	}

	// update connections list
	err = state.SAdd(eventProps.EventCtx, "connections", eventProps.Address.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func (plugin *MSyncPlugin) onDisconnect(eventProps events.EventProps) error {
	state, err := redis.GetRedisClient()
	if err != nil {
		return err
	}

	// remove client info
	state.Del(eventProps.EventCtx, "connection:"+eventProps.Address.String())
	// update connections list
	state.SRem(eventProps.EventCtx, "connections", eventProps.Address.String())

	return nil
}

func (plugin *MSyncPlugin) updateConnection(eventProps events.EventProps) error {
	// get redis client
	state, err := redis.GetRedisClient()
	if err != nil {
		return err
	}

	// get current state
	connectionJsonString, err := state.Get(eventProps.EventCtx, "connection:"+eventProps.Address.String()).Result()
	if err != nil {
		return err
	}

	// parse
	var connection MSyncConnection
	err = json.Unmarshal([]byte(connectionJsonString), &connection)
	if err != nil {
		return err
	}

	// update
	connection.UpdatedAt = time.Now().UnixMilli()
	// convert
	connectionJsonStringUpdate, err := json.Marshal(connection)
	if err != nil {
		return err
	}

	// save (update state in redis)
	err = state.Set(eventProps.EventCtx, "connection:"+eventProps.Address.String(), connectionJsonStringUpdate, connectionExpirationTime).Err()
	if err != nil {
		return err
	}

	return nil
}

var PluginInstance MSyncPlugin
