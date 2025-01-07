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
	"moony/plugins/auth/validator"
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
		// parse & validate input
		input, _, err := validator.ValidateCreateInput(data)
		if err != nil {
			events.SendError(plugin.config, "create", err.Error(), eventProps)
			return
		}

		// create user
		result, err := plugin.create(ctx, input)

		// return error if not created
		if err != nil {
			events.SendError(plugin.config, "create", err.Error(), eventProps)
			return
		}

		// return result
		events.Send(plugin.config, "create", []any{result.ID}, eventProps)
	})

	events.Create(plugin.config, "login", func(data []any, eventProps events.EventProps) {
		// parse & validate input
		input, _, err := validator.ValidateLoginInput(data)
		if err != nil {
			events.SendError(plugin.config, "login", err.Error(), eventProps)
			return
		}

		// authorize user
		result, err := plugin.login(ctx, input)

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

func (plugin *AuthPlugin) login(ctx context.Context, input *validator.AuthLoginInput) (*sqlc.User, error) {
	qc, err := queries_client.GetQueriesClient()
	if err != nil {
		return nil, err
	}

	commonError := errors.New("invalid_username_or_password")

	uniqueUser, err := qc.UsersFindUnique(ctx, input.Username)
	if err != nil {
		return nil, commonError
	}

	passwordValidate, err := crypto.HashValidate(input.Password, uniqueUser.Password)
	if err != nil {
		return nil, commonError
	}

	if passwordValidate {
		log.Println("uniqueUser", uniqueUser)
		return &uniqueUser, nil
	}

	return nil, commonError
}

func (plugin *AuthPlugin) create(ctx context.Context, input *validator.AuthCreateInput) (*sqlc.User, error) {
	qc, err := queries_client.GetQueriesClient()
	if err != nil {
		return nil, err
	}

	passwordHash, err := crypto.HashCreate(input.Password)
	if err != nil {
		return nil, err
	}

	insertedUser, err := qc.UsersCreate(ctx, sqlc.UsersCreateParams{
		Username: input.Username,
		Password: passwordHash,
		Email:    input.Email,
	})

	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}

var PluginInstance AuthPlugin
