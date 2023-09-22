/*
 license x
*/

package internal

import (
	"fmt"
	"net/http"
)

// ParsingError indicates that an error has occurred when parsing request parameters
type ParsingError struct {
	Err error
}

func (e *ParsingError) Unwrap() error {
	return e.Err
}

func (e *ParsingError) Error() string {
	return e.Err.Error()
}

// MethodNotAllowedError indicates that an error has occurred when parsing request parameters
type MethodNotAllowedError struct {
	Err error
}

func (e *MethodNotAllowedError) Error() string {
	return e.Err.Error()
}

// StorageError indicates that an error has occurred when parsing request parameters
type StorageError struct {
	Err error
}

func (e *StorageError) Error() string {
	return e.Err.Error()
}

// RequiredError indicates that an error has occurred when parsing request parameters
type RequiredError struct {
	Field string
}

func (e *RequiredError) Error() string {
	return fmt.Sprintf("required field '%s' is zero value.", e.Field)
}

// NotFoundError indicates that an error has occurred when parsing request parameters
type NotFoundError struct {
	Err error
}

func (e *NotFoundError) Error() string {
	return e.Err.Error()
}

// ErrorHandler defines the required method for handling error. You may implement it and inject this into a controller if
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// DefaultErrorHandler defines the default logic on how to handle errors from the controller. Any errors from parsing
func DefaultErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {
	var statusCode int

	switch err.(type) {
	case *ParsingError:
		statusCode = http.StatusBadRequest
	case *RequiredError:
		statusCode = http.StatusUnprocessableEntity
	case *MethodNotAllowedError:
		statusCode = http.StatusMethodNotAllowed
	case *NotFoundError:
		statusCode = http.StatusNotFound
	case *StorageError:
		statusCode = http.StatusBadGateway
	default:
		statusCode = http.StatusInternalServerError
	}

	_ = EncodeJSONResponse(func(i int) *int { return &i }(statusCode), w)
}

// EncodeJSONResponse uses the json encoder to write an interface to the http response with an optional status code
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
