/*
 license x
*/

package internal

import (
	"github.com/nelsonstr/o801/internal/log"
	"net/http"
)

// A Route defines the parameters for an api endpoint.
type Route struct {
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is a map of defined api endpoints.
type Routes map[string]Route

// Router defines the required methods for retrieving api routes
type Router interface {
	Routes() Routes
}

// NewRouter creates a new router for any number of api routers.
func NewRouter(routers ...Router) *http.ServeMux {
	mux := http.NewServeMux()

	for _, api := range routers {
		for name, route := range api.Routes() {
			var handler http.Handler
			handler = route.HandlerFunc

			mux.Handle(route.Pattern, log.Logger(handler, name))
		}
	}

	return mux
}

// EncodeJSONResponse uses the json encoder to write an interface to the http response with an optional status code
func EncodeJSONResponse(status *int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != nil {
		w.WriteHeader(*status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	return nil
}
