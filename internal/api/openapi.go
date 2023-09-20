package api

import "net/http"

func (r router) OpenAPI() router {
	r.Handle("/", http.FileServer(http.Dir("api")))

	return r
}
