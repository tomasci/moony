package validator

import (
	"moony/moony/core/mvalidator"
)

type AuthLoginInput struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthCreateInput struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
}

func ValidateLoginInput(data []any) (*AuthLoginInput, []string, error) {
	input := AuthLoginInput{
		Username: data[0].(string),
		Password: data[1].(string),
	}

	return mvalidator.Validate[AuthLoginInput](input)
}

func ValidateCreateInput(data []any) (*AuthCreateInput, []string, error) {
	input := AuthCreateInput{
		Username: data[0].(string),
		Password: data[1].(string),
		Email:    data[2].(string),
	}

	return mvalidator.Validate[AuthCreateInput](input)
}
