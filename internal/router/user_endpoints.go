/*
 license x
*/

package router

import (
	"github.com/nelsonstr/o801/api"
	"net/http"
)

func (r *Router) UserEndpoints(s api.HandlerFuncAPI) {
	r.Endpoint(http.MethodGet, "/(\\d+)+", s.Get)
	r.Endpoint(http.MethodPost, "/", s.Create)
}
