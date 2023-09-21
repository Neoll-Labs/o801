package api

import (
	"net/http"
)

func (r *router) AuthEndpoints() {
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("login"))
	})
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("logout"))
	})

}
