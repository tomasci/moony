package main

type MSyncPlugin struct{}

func Init() error {
	return nil
}

var PluginInstance MSyncPlugin
