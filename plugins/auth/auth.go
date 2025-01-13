package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"moony/database/queries_client"
	"moony/database/redis"
	"moony/database/sqlc"
	"moony/moony/core/crypto"
	"moony/moony/core/dispatcher"
	"moony/moony/core/events"
	"moony/moony/core/plugins"
	"moony/plugins/auth/validator"
	"net"
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
	d := dispatcher.GetGlobalDispatcher()

	events.Create(plugin.config, "create", func(data []any, eventProps events.EventProps) {
		// parse & validate input
		input, _, err := validator.ValidateCreateInput(data)
		if err != nil {
			events.SendError(plugin.config, "create", err.Error(), eventProps)
			return
		}

		// create user
		result, err := plugin.create(eventProps, input)

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
		result, err := plugin.login(eventProps, input)

		// return error is username or password is incorrect or user doesn't exist
		if err != nil {
			events.SendError(plugin.config, "login", err.Error(), eventProps)
			return
		}

		// return result
		// todo: send token
		events.Send(plugin.config, "login", []any{result.ID}, eventProps)
	})

	d.RegisterEventHandler("msync_disconnect", func(ctx context.Context, conn *net.UDPConn, address *net.UDPAddr, data []any) {
		plugin.cleanup(ctx, address)
	})

	return nil
}

func (plugin *AuthPlugin) login(eventProps events.EventProps, input *validator.AuthLoginInput) (*sqlc.User, error) {
	// get redis client
	state, err := redis.GetRedisClient()
	if err != nil {
		return nil, err
	}

	qc, err := queries_client.GetQueriesClient()
	if err != nil {
		return nil, err
	}

	commonError := errors.New("invalid_username_or_password")

	uniqueUser, err := qc.UsersFindUnique(eventProps.EventCtx, input.Username)
	if err != nil {
		return nil, commonError
	}

	passwordValidate, err := crypto.HashValidate(input.Password, uniqueUser.Password)
	if err != nil {
		return nil, commonError
	}

	if passwordValidate {
		// cache user info
		var uniqueUserJson []byte
		uniqueUserJson, err = json.Marshal(uniqueUser)
		if err != nil {
			return nil, err
		}

		err = state.Set(eventProps.EventCtx, "user:"+eventProps.Address.String(), uniqueUserJson, 0).Err()
		if err != nil {
			return nil, err
		}

		return &uniqueUser, nil
	}

	return nil, commonError
}

func (plugin *AuthPlugin) create(eventProps events.EventProps, input *validator.AuthCreateInput) (*sqlc.User, error) {
	qc, err := queries_client.GetQueriesClient()
	if err != nil {
		return nil, err
	}

	passwordHash, err := crypto.HashCreate(input.Password)
	if err != nil {
		return nil, err
	}

	insertedUser, err := qc.UsersCreate(eventProps.EventCtx, sqlc.UsersCreateParams{
		Username: input.Username,
		Password: passwordHash,
		Email:    input.Email,
	})

	if err != nil {
		return nil, err
	}

	return &insertedUser, nil
}

func (plugin *AuthPlugin) cleanup(ctx context.Context, address *net.UDPAddr) {
	// get redis client
	state, err := redis.GetRedisClient()
	if err != nil {
		log.Println("error getting redis client")
		return
	}

	state.Del(ctx, "user:"+address.String())
}

var PluginInstance AuthPlugin
