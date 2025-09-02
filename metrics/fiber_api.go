package metrics

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of requests by route and method",
		},
		[]string{"route", "method"},
	)

	responseCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_responses_total",
			Help: "Total number of responses by route, method and status code",
		},
		[]string{"route", "method", "status_code"},
	)

	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_latency_seconds",
			Help:    "Histogram of API request latencies by route",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route", "method"},
	)
)

func NewFiberApiMetricsMiddleware() func() fiber.Handler {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(responseCounter)
	prometheus.MustRegister(latencyHistogram)

	return fiberAPIMetricsMiddleware
}

func fiberAPIMetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		method := c.Method()
		path := normalizePath(c.Path())

		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()

		requestCounter.WithLabelValues(path, method).Inc()

		statusCode := c.Response().StatusCode()
		if statusCode >= http.StatusInternalServerError || err != nil {
			responseCounter.WithLabelValues(path, method, "500").Inc()
		} else {
			responseCounter.WithLabelValues(path, method, strconv.Itoa(statusCode)).Inc()
		}

		latencyHistogram.WithLabelValues(path, method).Observe(duration)

		return err
	}
}

var (
	uuidRegex    = regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)
	numericRegex = regexp.MustCompile(`/\d+(/|$)`)
)

func normalizePath(path string) string {
	normalizedPath := uuidRegex.ReplaceAllString(path, "<uuid>")
	normalizedPath = numericRegex.ReplaceAllString(normalizedPath, "/<numeric_id>")

	return normalizedPath
}
