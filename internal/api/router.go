package api

import (
	"github.com/nelsonstr/o801/internal/log"
	"github.com/nelsonstr/o801/internal/monitoring"
	"net/http"
	"regexp"
	"strconv"
)

type Route struct {
	Method   string
	Handler  http.HandlerFunc
	Patterns *regexp.Regexp
}

type Router struct {
	Mux      *http.ServeMux
	routes   []Route
	prefix   string
	resource string
}

func NewRouter() *Router {
	return &Router{
		Mux:    http.NewServeMux(),
		routes: []Route{},
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if req.Method != route.Method {
			continue
		}

		if !route.Patterns.MatchString(req.URL.Path) {
			continue
		}

		route.Handler.ServeHTTP(w, req)
		return
	}

	http.NotFound(w, req)
}

func (r *Router) Version(v int) *Router {
	vPath := "/api/v" + strconv.Itoa(v)
	r.Mux.Handle(vPath+"/", log.Logger(http.StripPrefix(vPath, r.Mux)))
	r.prefix = vPath

	return r
}

func (r *Router) Resource(name string) *Router {
	r.Mux.Handle(name+"/", monitoring.PrometheusMiddleware(http.StripPrefix(name, r.Mux)))
	r.resource = name

	return r
}

func (r *Router) Endpoint(method, path string, h http.HandlerFunc) {
	compiled, err := regexp.Compile("^" + r.prefix + r.resource + path + "$")
	if err != nil {
		panic(err)
	}

	r.routes = append(r.routes, Route{

		Patterns: compiled,
		Method:   method,
		Handler:  h,
	})

}
