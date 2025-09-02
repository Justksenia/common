package middlewares

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

func Labeler(find RouteFinder) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route, ok := find(r.Method, r.URL)
			if !ok {
				h.ServeHTTP(w, r)
				return
			}

			attr := semconv.HTTPRouteKey.String(route.PathPattern())
			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(attr)

			labeler, _ := otelhttp.LabelerFromContext(r.Context())
			labeler.Add(attr)

			h.ServeHTTP(w, r)
		})
	}
}
