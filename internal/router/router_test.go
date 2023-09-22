package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter_ServeHTTP(t *testing.T) {
	// Create a new router
	router := NewRouter()

	// Define a test handler
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Add a test route
	router.Endpoint("GET", "/test", testHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRouter_ServeHTTP404(t *testing.T) {
	// Create a new router
	router := NewRouter()

	// Define a test handler
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Add a test route
	router.Endpoint("GET", "/test", testHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test404", nil)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
func TestRouter_Version(t *testing.T) {
	// Create a new router
	router := NewRouter()

	// Apply a version
	router.Version(1)

	// Check if the prefix has been set correctly
	assert.Equal(t, "/api/v1", router.prefix)
}

func TestRouter_Resource(t *testing.T) {
	// Create a new router
	router := NewRouter()

	// Apply a resource
	router.Resource("/users")

	// Check if the resource has been set correctly
	assert.Equal(t, "/users", router.resource)
}

func TestRouter_Endpoint(t *testing.T) {
	// Create a new router
	router := NewRouter()

	// Define a test handler
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Add an endpoint
	router.Endpoint("GET", "/test", testHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRouter_ServeHTTP_NotFound(t *testing.T) {
	// Create a new router
	router := NewRouter()

	// Create a test request for an undefined route
	req := httptest.NewRequest("GET", "/undefined", nil)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Check the response status code (should be 404)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRouter_ServeHTTP_WithParams(t *testing.T) {
	// Create a new router
	router := NewRouter()

	// Define a test handler with parameters
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value("params").([]string)
		assert.Equal(t, "123", params[1])
		w.WriteHeader(http.StatusOK)
	}

	// Add a route with a parameterized path
	router.Endpoint("GET", "/test/([0-9]+)", testHandler)

	req := httptest.NewRequest("GET", "/test/123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
