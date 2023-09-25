/*
 license x
*/

package routes

import (
	"github.com/nelsonstr/o801/api"
	sql "github.com/nelsonstr/o801/internal/repository"
	"github.com/nelsonstr/o801/internal/router"
	"github.com/nelsonstr/o801/internal/user/handlers"
	"github.com/nelsonstr/o801/internal/user/repository"
	"github.com/nelsonstr/o801/internal/user/service"
	"net/http"
)

func InitUserRoutes(dbc sql.DBInterface, r *router.Router) {
	repository := repository.NewUserRepository(dbc)
	service := service.NewUserService(repository)
	handlers := handlers.NewUserHandler(service)
	UserEndpoints(r.Resource("/users"), handlers)
}

func UserEndpoints(r *router.Router, s api.HandlerAPI) {
	r.Endpoint(http.MethodGet, "/(\\d+)+", s.Get)
	r.Endpoint(http.MethodPost, "/", s.Create)
}
