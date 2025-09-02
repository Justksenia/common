package fibermw

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	cmnlogger "gitlab.com/adstail/ts-common/logger"
	"gitlab.com/adstail/ts-common/tracer"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	internalErrorCode = 500
	clientErrorCode   = 400
	redirectionCode   = 300
)

func LoggingServer(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.StartSpan(c.UserContext(), string(c.Request().RequestURI()), trace.SpanKindServer)
		defer span.End()

		start := time.Now()
		newLog := logger.With(
			zap.String("remote_ip", c.IP()),
			zap.String("host", c.Hostname()),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.String("request", fmt.Sprintf("%s %s", c.Method(), c.OriginalURL())),
			zap.String("trace-id", span.SpanContext().TraceID().String()),
			zap.String("request-id", c.GetReqHeaders()[fiber.HeaderXRequestID]),
		)

		c.SetUserContext(cmnlogger.ToContext(ctx, newLog))

		res := c.Response()
		fields := []zapcore.Field{
			zap.String("time", time.Since(start).String()),
		}

		err := c.Next()

		if err != nil {
			fields = append(fields, zap.Error(err))
		}
		fields = append(fields, zap.Int("status", res.StatusCode()))
		baseLog := newLog.With(zap.String("Tag", "HTTP middleware"))

		switch status := res.StatusCode(); {
		case status >= internalErrorCode:
			baseLog.Error("Server error", fields...)
		case status >= clientErrorCode:
			baseLog.Warn("Client error", fields...)
		case status >= redirectionCode:
			baseLog.Debug("Redirection", fields...)
		default:
			baseLog.Debug("Success", fields...)
		}

		return err
	}
}
