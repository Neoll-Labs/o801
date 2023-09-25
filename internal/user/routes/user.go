/*
 license x
*/

package routes

import (
	"net/http"
	"time"

	"github.com/nelsonstr/o801/internal/cache"
	"github.com/nelsonstr/o801/internal/interfaces"
	"github.com/nelsonstr/o801/internal/model"
	repo "github.com/nelsonstr/o801/internal/repository"
	"github.com/nelsonstr/o801/internal/router"
	"github.com/nelsonstr/o801/internal/user/handlers"
	"github.com/nelsonstr/o801/internal/user/repository"
	"github.com/nelsonstr/o801/internal/user/service"
)

func InitUserRoutes(dbc repo.DBInterface, r *router.Router) {
	rep := repository.NewUserRepository(dbc)

	const ttl = 5 * time.Second
	c := cache.NewCache[model.UserView](ttl)

	serv := service.NewUserService(rep, c)

	hand := handlers.NewUserHandler(serv)

	UserEndpoints(r.Resource("/users"), hand)
}

func UserEndpoints(r *router.Router, s interfaces.HandlerAPI) {
	r.Endpoint(http.MethodGet, "/(\\d+)+", s.Get)
	r.Endpoint(http.MethodPost, "/", s.Create)
}
