/*
 license x
*/

package router

import (
	"github.com/nelsonstr/o801/internal"
	"net/http"
)

func (r *Router) UserEndpoints(s internal.HandlerFuncAPI) {
	r.Endpoint(http.MethodGet, "/(\\d+)+", s.Get)
	r.Endpoint(http.MethodPost, "/", s.Create)
}
