package api

import (
	"github.com/nelsonstr/o801/internal/handlers"
	"net/http"
)

func (r *Router) UserEndpoints(s *handlers.Server) {
	r.Endpoint(http.MethodGet, "/(\\d+)+", s.GetUser)
	r.Endpoint(http.MethodPost, "/", s.CreateUser)
}
