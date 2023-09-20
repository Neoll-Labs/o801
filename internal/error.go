/*
 license x
*/

package internal

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	// ErrTypeAssertionError is thrown when type an interface does not match the asserted type
	ErrTypeAssertionError = errors.New("unable to assert type")
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

// MethodNotAllowedError indicates that an error has occurred when parsing request parameters
type BadRequestError struct {
	Err error
}

func (e *BadRequestError) Error() string {
	return e.Err.Error()
}

// RequiredError indicates that an error has occurred when parsing request parameters
type RequiredError struct {
	Field string
}

func (e *RequiredError) Error() string {
	return fmt.Sprintf("required field '%s' is zero value.", e.Field)
}

// ErrorHandler defines the required method for handling error. You may implement it and inject this into a controller if
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error, result *ImplResponse)

// DefaultErrorHandler defines the default logic on how to handle errors from the controller. Any errors from parsing
func DefaultErrorHandler(w http.ResponseWriter, _ *http.Request, err error, result *ImplResponse) {
	if _, ok := err.(*ParsingError); ok {
		// Handle parsing errors
		log.Println(err.Error())
		_ = EncodeJSONResponse(func(i int) *int { return &i }(http.StatusBadRequest), w)
	} else if _, ok := err.(*RequiredError); ok {
		// Handle missing required errors
		log.Println(err.Error())
		_ = EncodeJSONResponse(func(i int) *int { return &i }(http.StatusUnprocessableEntity), w)
	} else if _, ok := err.(*MethodNotAllowedError); ok {
		// Handle method not allowed errors
		log.Println(err.Error())
		_ = EncodeJSONResponse(func(i int) *int { return &i }(http.StatusMethodNotAllowed), w)
	} else {
		// Handle all other errors
		log.Println(err.Error())
		_ = EncodeJSONResponse(&result.Code, w)
	}
}
