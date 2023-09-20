package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (r *router) Metrics() *router {
	r.Handle("/metrics/", promhttp.Handler())

	return r
}
