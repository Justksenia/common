package middlewares

import (
	"net/http"

	cmlogger "gitlab.com/adstail/ts-common/logger"
	"go.uber.org/zap"
)

// InjectLogger injects logger into request context.
func InjectLogger(lg *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCtx := r.Context()
			req := r.WithContext(cmlogger.ToContext(reqCtx, lg))
			next.ServeHTTP(w, req)
		})
	}
}
