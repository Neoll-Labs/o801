/*
 license x
*/

package internal

import (
	"errors"
	"net/http"
)

// ParsingError indicates that an error has occurred when parsing request parameters.
type ParsingError struct {
	Err error
}

func (e *ParsingError) Error() string {
	return e.Err.Error()
}

type MethodNotAllowedError struct {
	Err error
}

func (e *MethodNotAllowedError) Error() string {
	return e.Err.Error()
}

type StorageError struct {
	Err error
}

func (e *StorageError) Error() string {
	return e.Err.Error()
}

type NotFoundError struct {
	Err error
}

func (e *NotFoundError) Error() string {
	return e.Err.Error()
}

// ErrorHandler defines the required method for handling error. You may implement it and inject this into a controller.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// DefaultErrorHandler defines the default logic on how to handle errors from the controller.
func DefaultErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	var statusCode int

	var parsingError *ParsingError

	var notFoundError *NotFoundError

	var storageError *StorageError

	switch {
	case errors.As(err, &parsingError):
		statusCode = http.StatusBadRequest
	case errors.As(err, &notFoundError):
		statusCode = http.StatusNotFound
	case errors.As(err, &storageError):
		statusCode = http.StatusBadGateway
	default:
		statusCode = http.StatusInternalServerError
	}

	_ = EncodeJSONResponse(func(i int) *int { return &i }(statusCode), w)
}

// EncodeJSONResponse uses the json encoder to write an interface to the http response with an optional status code.
func EncodeJSONResponse(status *int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if status != nil {
		code := *status
		w.WriteHeader(code)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	return nil
}
