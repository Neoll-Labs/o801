package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter_ServeHTTP(t *testing.T) {
	router := NewRouter()

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	router.Endpoint("GET", "/test", testHandler)

	req := httptest.NewRequest("GET", "/test", nil)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRouter_ServeHTTP404(t *testing.T) {
	router := NewRouter()

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	router.Endpoint("GET", "/test", testHandler)
	req := httptest.NewRequest("GET", "/test404", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
func TestRouter_Version(t *testing.T) {
	router := NewRouter()
	router.Version(1)
	assert.Equal(t, "/api/v1", router.prefix)
}

func TestRouter_Resource(t *testing.T) {
	router := NewRouter()
	router.Resource("/users")
	assert.Equal(t, "/users", router.resource)
}

func TestRouter_Endpoint(t *testing.T) {
	router := NewRouter()
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	router.Endpoint("GET", "/test", testHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRouter_ServeHTTP_NotFound(t *testing.T) {
	router := NewRouter()
	req := httptest.NewRequest("GET", "/undefined", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRouter_ServeHTTP_WithParams(t *testing.T) {
	router := NewRouter()
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value("params").([]string)
		assert.Equal(t, "123", params[1])
		w.WriteHeader(http.StatusOK)
	}
	router.Endpoint("GET", "/test/([0-9]+)", testHandler)
	req := httptest.NewRequest("GET", "/test/123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
