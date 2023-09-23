/*
 license x
*/

package router

import (
	"net/http"

	"github.com/nelsonstr/o801/api"
)

func (r *Router) UserEndpoints(s api.HandlerFuncAPI) {
	r.Endpoint(http.MethodGet, "/(\\d+)+", s.Get)
	r.Endpoint(http.MethodPost, "/", s.Create)
}
