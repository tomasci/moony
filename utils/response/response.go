package response

import (
	"encoding/json"
)

type Response[DataType interface{}] struct {
	// Status is integer code, will probably look similar to codes in HTTP.
	// Basic statuses are: 500 - internal error, 200 - ok
	Status int `json:"status"`
	// Error is just boolean true/false
	Error *bool `json:"error,omitempty"`
	// Message is description of what happened (only for errors, do not use it for anything else)
	Message *string `json:"message,omitempty"`
	// Plugin is for plugin name
	Plugin string `json:"plugin"`
	// Method is for plugin method
	Method string `json:"method"`
	// Data is for any data
	Data DataType `json:"data,omitempty"`
}

func Success[DataType interface{}](plugin string, method string, data DataType) ([]byte, error) {
	response := Response[DataType]{
		Status:  200,
		Error:   nil,
		Message: nil,
		Plugin:  plugin,
		Method:  method,
		Data:    data,
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func Error[DataType interface{}](status int, plugin string, method string, data DataType, err error) ([]byte, error) {
	error := true
	errorMsg := err.Error()

	response := Response[any]{
		Status:  status,
		Error:   &error,
		Message: &errorMsg,
		Plugin:  plugin,
		Method:  method,
		Data:    data,
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}
