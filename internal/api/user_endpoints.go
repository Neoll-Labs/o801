package api

import (
	"net/http"
)

func (r *router) UserEndpoints() {

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("john doe"))
	})

}
