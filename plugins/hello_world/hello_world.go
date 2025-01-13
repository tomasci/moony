package main

import (
	"context"
	"github.com/google/uuid"
	"log"
	"moony/database/redis"
	"moony/moony/bridges/godot"
	"moony/moony/bridges/godot/materials"
	"moony/moony/core/dispatcher"
	"moony/moony/core/events"
	"moony/moony/core/plugins"
	"net"
	"strings"
	"unicode"
)

type HelloWorldPlugin struct {
	config plugins.PluginConfig
}

func init() {
	// you can use this function if you need to do something even before Init method called
	log.Println("Hello World Plugin init() func")
	//time.Sleep(1 * time.Second) // simulate some hard work during initialization...
}

func (plugin *HelloWorldPlugin) Init(ctx context.Context, config plugins.PluginConfig) error {
	log.Println("Hello World Plugin Init() func")
	plugin.config = config
	d := dispatcher.GetGlobalDispatcher()

	// with "events" introduction, you don't need to use "raw" dispatcher anymore
	// but, if you want to subscribe to some internal events (like events from server)
	// you still can do it with d and GetGlobalDispatcher
	d.RegisterEventHandler("OnServerStarted", func(eventCtx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any) {
		log.Printf("%v: %v", plugin.config.Name, data)
	})

	// for everything else use events
	events.Create(plugin.config, "sum", func(data []any, eventProps events.EventProps) {
		// parse input
		a, aOk := data[0].(int)
		b, bOk := data[1].(int)

		// send error
		if !aOk || !bOk {
			events.SendError(plugin.config, "sum", "invalid_input_data", eventProps)
			return
		}

		// calculate
		result := plugin.sum(a, b)

		// return result
		events.Send(plugin.config, "sum", []any{result}, eventProps)
	})

	events.Create(plugin.config, "capitalize", func(data []any, eventProps events.EventProps) {
		// parse input
		str, ok := data[0].(string)

		// send error
		if !ok {
			events.SendError(plugin.config, "capitalize", "invalid_input_data", eventProps)
			return
		}

		// transform string
		result := plugin.capitalize(str)

		// return result
		//events.Send(plugin.config, "capitalize", []any{result}, eventProps)
		// return result to everyone
		events.Broadcast(plugin.config, "capitalize", []any{result}, eventProps)
	})

	events.Create(plugin.config, "spawn_object", func(data []any, eventProps events.EventProps) {
		result, err := plugin.spawnObject(eventProps)
		if err != nil {
			log.Println("spawn_object error:", err)
			return
		}

		events.Broadcast(plugin.config, "spawn_object", result, eventProps)
	})

	return nil
}

func (plugin *HelloWorldPlugin) sum(a, b int) int {
	return a + b
}

func (plugin *HelloWorldPlugin) capitalize(str string) string {
	// some capitalize code from internet
	words := strings.Fields(str)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = string(unicode.ToUpper(rune(word[0]))) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

func (plugin *HelloWorldPlugin) spawnObject(eventProps events.EventProps) ([]any, error) {
	// get redis client
	state, err := redis.GetRedisClient()
	if err != nil {
		return nil, err
	}

	// generate object id
	objectId := uuid.New()

	// track objects
	// just so you know how much object spawned, maybe for cleanup etc.
	// maybe for statistics, who knows
	err = state.SAdd(eventProps.EventCtx, "world_objects", objectId.String()).Err()
	if err != nil {
		return nil, err
	}

	// todo: first need to create object
	// only then - spawn
	// also, need to set object id
	object := godot.SpawnObject(godot.Position{
		X: godot.RandomFloat(-3.50, 3.50),
		Y: godot.RandomFloat(-3.50, 3.50),
		Z: godot.RandomFloat(-3.50, 3.50),
	}, godot.Object{
		ID: objectId,
	}, materials.Random())

	return object, nil
}

var PluginInstance HelloWorldPlugin
