/*
 license x
*/

/*
 * O801 API
 *
 * Create and Get User
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
)

// A Route defines the parameters for an api endpoint
type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is a map of defined api endpoints
type Routes map[string]Route

// Router defines the required methods for retrieving api routes
type Router interface {
	Routes() Routes
}

// NewRouter creates a new router for any number of api routers
func NewRouter(routers ...Router) *http.ServeMux {
	mux := http.NewServeMux()

	for _, api := range routers {
		for name, route := range api.Routes() {
			var handler http.Handler
			handler = route.HandlerFunc

			mux.Handle(route.Pattern, Logger(handler, name))
		}
	}

	return mux
}

// EncodeJSONResponse uses the json encoder to write an interface to the http response with an optional status code
func EncodeJSONResponse(i interface{}, status *int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != nil {
		w.WriteHeader(*status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if i != nil {
		return json.NewEncoder(w).Encode(i)
	}

	return nil
}

type Number interface {
	~int32 | ~int64 | ~float32 | ~float64
}

type ParseString[T Number | string | bool] func(v string) (T, error)

type Operation[T Number | string | bool] func(actual string) (T, bool, error)

type Constraint[T Number | string | bool] func(actual T) error
