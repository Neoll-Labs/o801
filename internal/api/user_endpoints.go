package api

import (
	"github.com/nelsonstr/o801/internal/server"
	"net/http"
)

func (r *Router) UserEndpoints(s *server.Server) {

	r.Endpoint(http.MethodGet, "/(\\d+)+", s.GetUser)
	r.Endpoint(http.MethodPost, "/", s.CreateUser)
}
