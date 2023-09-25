/*
 license x
*/

package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter_ServeHTTP(t *testing.T) {
	t.Parallel()
	router := NewRouter()

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	router.Endpoint(http.MethodGet, "/test", testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRouter_ServeHTTP404(t *testing.T) {
	t.Parallel()
	router := NewRouter()

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	router.Endpoint(http.MethodGet, "/test", testHandler)
	req := httptest.NewRequest(http.MethodGet, "/test404", http.NoBody)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRouter_Version(t *testing.T) {
	t.Parallel()
	router := NewRouter()
	router.Version(1)
	assert.Equal(t, "/api/v1", router.prefix)
}

func TestRouter_Resource(t *testing.T) {
	t.Parallel()
	router := NewRouter()
	router.Resource("/users")
	assert.Equal(t, "/users", router.resource)
}

func TestRouter_Endpoint(t *testing.T) {
	t.Parallel()
	router := NewRouter()
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	router.Endpoint(http.MethodGet, "/test", testHandler)
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRouter_MethodNotMatch(t *testing.T) {
	t.Parallel()
	router := NewRouter()
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	router.Endpoint(http.MethodGet, "/test", testHandler)
	req := httptest.NewRequest(http.MethodPost, "/test", http.NoBody)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestRouter_ServeHTTP_WithParams(t *testing.T) {
	t.Parallel()
	router := NewRouter()
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(ParametersName).([]string)
		assert.Equal(t, "123", params[1])
		w.WriteHeader(http.StatusOK)
	}
	router.Endpoint(http.MethodGet, "/test/([0-9]+)", testHandler)
	req := httptest.NewRequest(http.MethodGet, "/test/123", http.NoBody)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
