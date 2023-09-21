package api

import "net/http"

func (r router) OpenAPI() {
	r.Handle("/", http.FileServer(http.Dir("api")))

}
