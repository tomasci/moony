package main

import (
	"context"
	"log"
	"moony/moony/core/event_dispatcher"
	"moony/moony/core/plugins"
	"moony/moony/utils/response"
	"net"
	"strings"
	"time"
	"unicode"
)

type HelloWorldPlugin struct {
	config plugins.PluginConfig
}

func init() {
	// you can use this function if you need to do something even before Init method called
	log.Println("Hello World Plugin init() func")
	time.Sleep(1 * time.Second) // simulate some hard work during initialization...
}

func (plugin *HelloWorldPlugin) Init(ctx context.Context, dispatcher *event_dispatcher.EventDispatcher, config plugins.PluginConfig) error {
	log.Println("Hello World Plugin Init() func")
	plugin.config = config

	// this is how you usually register event handlers
	dispatcher.RegisterEventHandler("OnServerStarted", func(eventCtx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any) {
		log.Printf("%v: %v", plugin.config.Name, data)
	})

	// this is example, where you use plugin name as namespace and dot to separate plugin name and method name
	dispatcher.RegisterEventHandler(plugin.config.Name+".sum", func(eventCtx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any) {
		a, aOk := data[0].(int)
		b, bOk := data[1].(int)
		if !aOk || !bOk {
			log.Println(plugin.config.Name+".sum.result", "invalid input data")
			return
		}

		result := plugin.sum(a, b)

		// and when you want to return result, you call dispatch with the same event name, but .result suffix
		dispatcher.Dispatch(plugin.config.Name+".sum.result", eventCtx, conn, address, []any{result})
	})

	dispatcher.RegisterEventHandler(plugin.config.Name+".capitalize", func(eventCtx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any) {
		str, ok := data[0].(string)
		if !ok {
			log.Println(plugin.config.Name+".capitalize.result", "invalid input data")
			return
		}

		result := plugin.capitalize(str)
		if conn != nil && address != nil {
			response.SendResponse[any](conn, address, plugin.config.Name, "capitalize", result, nil)
		}
		dispatcher.Dispatch(plugin.config.Name+".capitalize.result", eventCtx, conn, address, []any{result})
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

var PluginInstance HelloWorldPlugin
