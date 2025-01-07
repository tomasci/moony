package mvalidator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"sync"
)

var (
	once     sync.Once
	validate *validator.Validate
)

func InitializeValidator() {
	once.Do(func() {
		log.Println("Initializing validator")
		validate = validator.New(validator.WithRequiredStructEnabled())
	})
}

func GetValidator() *validator.Validate {
	if validate == nil {
		InitializeValidator()
	}

	return validate
}

func Validate[InputType interface{}](input InputType) (*InputType, []string, error) {
	validate := GetValidator()
	err := validate.Struct(input)

	errorFields := make([]string, 0)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			structField, _ := reflect.TypeOf(input).FieldByName(err.StructField())
			jsonTag := structField.Tag.Get("json")
			errorFields = append(errorFields, jsonTag)
		}

		log.Println(errorFields)
		log.Println("validation_failed", err.Error())

		return nil, errorFields, errors.New("validation_failed")
	}

	return &input, errorFields, nil
}
