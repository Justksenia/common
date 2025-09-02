package middlewares

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Metrics wraps TracerProvider and MeterProvider.
type Metrics interface {
	TracerProvider() trace.TracerProvider
	MeterProvider() metric.MeterProvider
	TextMapPropagator() propagation.TextMapPropagator
}

// Instrument setups otelhttp.
func Instrument(serviceName string, find RouteFinder, m Metrics) Middleware {
	return func(h http.Handler) http.Handler {
		return otelhttp.NewHandler(h, "",
			otelhttp.WithPropagators(m.TextMapPropagator()),
			otelhttp.WithTracerProvider(m.TracerProvider()),
			otelhttp.WithMeterProvider(m.MeterProvider()),
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			otelhttp.WithServerName(serviceName),
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				op, ok := find(r.Method, r.URL)
				if ok {
					return serviceName + "." + op.OperationID()
				}
				return operation
			}),
		)
	}
}
