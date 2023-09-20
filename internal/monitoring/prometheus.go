/*
 license x
*/

package monitoring

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalRequests = prometheus.NewCounterVec( //nolint:gochecknoglobals // prometheus stats
		prometheus.CounterOpts{
			Namespace:   "nelsonstr",
			Subsystem:   "demo",
			Name:        "http_requests_total",
			Help:        "Number of get requests.",
			ConstLabels: nil,
		},
		[]string{"path"},
	)

	responseStatus = prometheus.NewCounterVec( //nolint:gochecknoglobals // prometheus stats
		prometheus.CounterOpts{
			Namespace:   "nelsonstr",
			Subsystem:   "demo",
			Name:        "response_status_total",
			Help:        "Status of HTTP response",
			ConstLabels: nil,
		},
		[]string{"status"},
	)

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{ //nolint:gochecknoglobals,exhaustruct // prometheus stats
		Namespace: "nelsonstr",
		Subsystem: "demo",
		Name:      "http_response_time_seconds",
		Help:      "Duration of HTTP requests.",
	}, []string{"path"})

	upTime = promauto.NewCounter(prometheus.CounterOpts{ //nolint:gochecknoglobals,exhaustruct // prometheus stats
		Namespace: "nelsonstr",
		Subsystem: "demo",
		Name:      "up_time_total",
		Help:      "The up time in sec",
	})
)

func init() { //nolint:gochecknoinits // init prometheus
	recordMetrics()

	_ = prometheus.Register(totalRequests)
	_ = prometheus.Register(responseStatus)
	_ = prometheus.Register(httpDuration)
}

func recordMetrics() {
	go func() {
		for {
			upTime.Inc()
			time.Sleep(1 * time.Second)
		}
	}()
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode
		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()
		timer.ObserveDuration()
	})
}
