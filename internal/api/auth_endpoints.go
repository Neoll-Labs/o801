package api

import (
	"net/http"
)

func (r *Router) AuthEndpoints() {
	r.Endpoint(http.MethodPost, "/login", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("login"))
	})

	r.Endpoint(http.MethodGet, "/logout", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("logout"))
	})

}
