package restprometheus

import (
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dflBuckets = []float64{300, 1200, 5000}
)

const (
	reqsName    = "json_rest_requests_total"
	latencyName = "json_rest_request_duration_milliseconds"
)

// Middleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type PromMiddleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
	name    string
}

// NewMiddleware returns a new prometheus Middleware handler.
func (mw *PromMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	start := time.Now()
	mw.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": mw.name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(mw.reqs)

	mw.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": mw.name},
		Buckets:     dflBuckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(mw.latency)

	return func(w rest.ResponseWriter, r *rest.Request) {
		h(w, r)

		if r.Env["STATUS_CODE"] == nil {
			log.Fatal("StatusMiddleware: Env[\"STATUS_CODE\"] is nil, " +
				"RecorderMiddleware may not be in the wrapped Middlewares.")
		}
		statusCode := r.Env["STATUS_CODE"].(int)

		mw.reqs.WithLabelValues(http.StatusText(statusCode), r.Method, r.URL.Path).Inc()
		mw.latency.WithLabelValues(http.StatusText(statusCode), r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}

}
