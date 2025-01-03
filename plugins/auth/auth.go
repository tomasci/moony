package main

import (
	"context"
	"moony/database/queries_client"
	"moony/database/sqlc"
	"moony/moony/core/events"
	"moony/moony/core/plugins"
)

type AuthPlugin struct {
	config plugins.PluginConfig
}

func init() {
	return
}

// this plugin looks like hell, but I must create plugins to see where the problems are
// and to understand how to make it better

func (plugin *AuthPlugin) Init(ctx context.Context, config plugins.PluginConfig) error {
	plugin.config = config

	events.Create(plugin.config, "create", func(data []any, eventProps events.EventProps) {
		// parse input
		username, uOk := data[0].(string)
		password, pOk := data[1].(string)
		email, eOk := data[2].(string)

		// validate
		if !uOk || !pOk || !eOk {
			events.SendError(plugin.config, "create", "invalid_input_data", eventProps)
		}

		// create user
		result, err := plugin.create(ctx, username, password, email)

		// return error if not created
		if err != nil {
			events.SendError(plugin.config, "create", err.Error(), eventProps)
		}

		// return result
		events.Send(plugin.config, "create", []any{result}, eventProps)
	})

	return nil
}

func (plugin *AuthPlugin) login() {
	return
}

func (plugin *AuthPlugin) create(ctx context.Context, username string, password string, email string) (*sqlc.User, error) {
	qc, err := queries_client.GetQueriesClient()
	if err != nil {
		return nil, err
	}

	insertedUser, err := qc.UsersCreate(ctx, sqlc.UsersCreateParams{
		Username: username,
		Password: password,
		Email:    email,
	})

	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}

var PluginInstance AuthPlugin
