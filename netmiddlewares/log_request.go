package middlewares

import (
	"net/http"

	cmnlogger "github.com/Justksenia/common/logger"
	"go.uber.org/zap"
)

// LogRequests logs incoming requests using context logger.
func LogRequests(find RouteFinder) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			logger := cmnlogger.FromContext(ctx)
			var (
				opName = zap.Skip()
				opID   = zap.Skip()
			)

			if route, ok := find(r.Method, r.URL); ok {
				opName = zap.String("operation_name", route.Name())
				opID = zap.String("operation_id", route.OperationID())
			}

			logger.Info("Got request", zap.String("method", r.Method), zap.Stringer("url", r.URL), opID, opName)
			next.ServeHTTP(w, r)
		})
	}
}
