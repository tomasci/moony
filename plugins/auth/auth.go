package main

import (
	"context"
	"errors"
	"log"
	"moony/database/queries_client"
	"moony/database/sqlc"
	"moony/moony/core/crypto"
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
			return
		}

		// create user
		result, err := plugin.create(ctx, username, password, email)

		// return error if not created
		if err != nil {
			events.SendError(plugin.config, "create", err.Error(), eventProps)
			return
		}

		// return result
		events.Send(plugin.config, "create", []any{result.ID}, eventProps)
	})

	events.Create(plugin.config, "login", func(data []any, eventProps events.EventProps) {
		// parse input
		username, uOk := data[0].(string)
		password, pOk := data[1].(string)

		// validate
		if !uOk || !pOk {
			events.SendError(plugin.config, "login", "invalid_input_data", eventProps)
			return
		}

		// authorize user
		result, err := plugin.login(ctx, username, password)

		// return error is username or password is incorrect or user doesn't exist
		if err != nil {
			events.SendError(plugin.config, "login", err.Error(), eventProps)
			return
		}

		// return result
		// todo: send token
		events.Send(plugin.config, "login", []any{result.ID}, eventProps)
	})

	return nil
}

func (plugin *AuthPlugin) login(ctx context.Context, username string, password string) (*sqlc.User, error) {
	qc, err := queries_client.GetQueriesClient()
	if err != nil {
		return nil, err
	}

	commonError := errors.New("invalid_username_or_password")

	uniqueUser, err := qc.UsersFindUnique(ctx, username)
	if err != nil {
		return nil, commonError
	}

	passwordValidate, err := crypto.HashValidate(password, uniqueUser.Password)
	if err != nil {
		return nil, commonError
	}

	if passwordValidate {
		log.Println("uniqueUser", uniqueUser)
		return &uniqueUser, nil
	}

	return nil, commonError
}

func (plugin *AuthPlugin) create(ctx context.Context, username string, password string, email string) (*sqlc.User, error) {
	qc, err := queries_client.GetQueriesClient()
	if err != nil {
		return nil, err
	}

	passwordHash, err := crypto.HashCreate(password)
	if err != nil {
		return nil, err
	}

	insertedUser, err := qc.UsersCreate(ctx, sqlc.UsersCreateParams{
		Username: username,
		Password: passwordHash,
		Email:    email,
	})

	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}

var PluginInstance AuthPlugin
