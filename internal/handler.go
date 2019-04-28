package internal

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	// requestCount is a simple counter that increments with each HTTP request
	// it is auto registered to the default metrics registry
	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "route", "status_code"})

	// duration is a histogram that is used to measure the response time (in seconds)
	// of http requests. The durations are bucketed into sane defaults
	// it is auto registered to the default metrics registry
	duration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route", "status_code"})
)

// NewHandler instantiates an http.Handler that has the relevant proxy and
// metrics route configured
func NewHandler(logger *logrus.Logger, targetAddr string) http.Handler {
	r := mux.NewRouter()
	h := handler{
		Addr:   targetAddr,
		Logger: logger,
	}

	r.Handle("/metrics", promhttp.Handler())
	r.PathPrefix("/").Handler(measure(h))

	return r
}

type handler struct {
	Addr   string
	Logger *logrus.Logger
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	targetURI := req.RequestURI
	uri := h.Addr + targetURI
	target, _ := url.Parse(uri)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", target.Host)
		req.URL.Scheme = "http"
		req.URL.Host = target.Host
	}

	h.Logger.Info(target)

	proxy := &httputil.ReverseProxy{
		Director: director,
	}

	proxy.ServeHTTP(w, req)
}

func measure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)

		route := r.RequestURI

		lvs := []string{r.Method, route, strconv.Itoa(m.Code)}

		requestCount.WithLabelValues(lvs...).Add(1)
		duration.WithLabelValues(lvs...).Observe(m.Duration.Seconds())
	})
}
