package mvalidator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"sync"
)

var (
	// to ensure instance created only once
	once sync.Once
	// validator instance
	validate *validator.Validate
)

// InitializeValidator – create validator instance
func InitializeValidator() {
	once.Do(func() {
		log.Println("Initializing validator")
		validate = validator.New(validator.WithRequiredStructEnabled())
	})
}

// GetValidator – can be used to get validator globally
func GetValidator() *validator.Validate {
	// check if instance exists, create if not
	if validate == nil {
		InitializeValidator()
	}

	return validate
}

func Validate[InputType interface{}](input InputType) (*InputType, []string, error) {
	// get validator
	validate := GetValidator()
	// validate with go-playground validator
	err := validate.Struct(input)

	// array for fields with errors
	errorFields := make([]string, 0)

	if err != nil {
		// go through all errors
		for _, err := range err.(validator.ValidationErrors) {
			// search field in struct
			structField, _ := reflect.TypeOf(input).FieldByName(err.StructField())
			// get it's json name
			jsonTag := structField.Tag.Get("json")
			// collect field names
			errorFields = append(errorFields, jsonTag)
		}

		//log.Println(errorFields)
		log.Println("validation_failed", err.Error())

		// return error
		return nil, errorFields, errors.New("validation_failed")
	}

	// return result with input
	return &input, errorFields, nil
}
