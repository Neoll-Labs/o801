package api

import (
	"net/http"
	"strconv"

	"github.com/nelsonstr/o801/internal/log"
	"github.com/nelsonstr/o801/internal/monitoring"
)

type router struct {
	*http.ServeMux
}

func NewRouter() *router {
	return &router{
		http.NewServeMux(),
	}
}

func (r *router) Prefix(path string) *router {
	mux := NewRouter()
	var handler http.Handler = mux

	r.Handle(path+"/", http.StripPrefix(path, handler))

	return mux
}

func (r *router) Resources(s string) *router {
	mux := NewRouter()
	var handler http.Handler = mux

	r.Handle(s+"/", monitoring.PrometheusMiddleware(http.StripPrefix(s, handler)))

	return mux
}

func (r *router) ApiVersion(v int) *router {
	mux := NewRouter()
	var handler http.Handler = mux

	vPath := "/api/v" + strconv.Itoa(v)
	r.Handle(vPath+"/", log.Logger(http.StripPrefix(vPath, handler)))

	return mux
}
