package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func (r *Router) Metrics() {
	r.Endpoint(http.MethodGet, "/metrics/", promhttp.Handler().ServeHTTP)

}
