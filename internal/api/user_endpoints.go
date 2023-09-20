package api

import (
	"net/http"
)

func (r *router) UserEndpoints() *router {

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("john doe"))
	})

	return r
}
