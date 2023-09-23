/*
 license x
*/

package internal

import (
	"errors"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestErrorTypes(t *testing.T) {
	tests := []struct {
		err      error
		expected string
	}{
		{&ParsingError{Err: errors.New("parsing error")}, "parsing error"},
		{&MethodNotAllowedError{Err: errors.New("method not allowed")}, "method not allowed"},
		{&StorageError{Err: errors.New("storage error")}, "storage error"},
		{&NotFoundError{Err: errors.New("not found error")}, "not found error"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, tt.err.Error())
			}
		})
	}
}

func TestDefaultErrorHandler(t *testing.T) {
	tests := []struct {
		err           error
		expectedCode  int
		expectedError string
	}{
		{&ParsingError{Err: errors.New("parsing error")}, http.StatusBadRequest, "parsing error"},
		{&NotFoundError{Err: errors.New("not found error")}, http.StatusNotFound, "not found error"},
		{&StorageError{Err: errors.New("storage error")}, http.StatusBadGateway, "storage error"},
		{errors.New("unknown error"), http.StatusInternalServerError, ""},
	}

	for _, tt := range tests {
		t.Run(tt.expectedError, func(t *testing.T) {
			w := httptest.NewRecorder()
			DefaultErrorHandler(w, nil, tt.err)
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestEncodeJSONResponse_WithStatus(t *testing.T) {
	w := httptest.NewRecorder()
	status := http.StatusCreated
	err := EncodeJSONResponse(&status, w)

	assert.NoError(t, err)
	assert.Equal(t, status, w.Code)

	contentType := w.Header().Get("Content-Type")
	assert.Equal(t, "application/json; charset=UTF-8", contentType)
}

func TestEncodeJSONResponse_WithoutStatus(t *testing.T) {
	w := httptest.NewRecorder()
	err := EncodeJSONResponse(nil, w)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	contentType := w.Header().Get("Content-Type")
	assert.Equal(t, contentType, "application/json; charset=UTF-8")
}
