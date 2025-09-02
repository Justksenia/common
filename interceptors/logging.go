package interceptors

import (
	"context"

	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	cmnlogger "gitlab.com/adstail/ts-common/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func UnaryLoggingInterceptor(logger *zap.Logger, opts ...grpcZap.Option) grpc.UnaryServerInterceptor {
	if logger == nil {
		var err error
		if logger, err = cmnlogger.New(); err != nil {
			panic(err)
		}
	}
	opts = append([]grpcZap.Option{defaultFormatter(), defaultDecider(logger.Level())}, opts...)
	return grpcZap.UnaryServerInterceptor(logger, opts...)
}

func StreamLoggingInterceptor(logger *zap.Logger, opts ...grpcZap.Option) grpc.StreamServerInterceptor {
	if logger == nil {
		var err error
		if logger, err = cmnlogger.New(); err != nil {
			panic(err)
		}
	}
	opts = append([]grpcZap.Option{defaultFormatter(), defaultDecider(logger.Level())}, opts...)
	return grpcZap.StreamServerInterceptor(logger, opts...)
}

func defaultDecider(ll zapcore.Level) grpcZap.Option {
	decider := func(_ string, err error) bool {
		return true
	}

	if ll > zapcore.WarnLevel {
		decider = func(_ string, err error) bool {
			return err != nil
		}
	}

	return grpcZap.WithDecider(decider)
}

func defaultFormatter() grpcZap.Option {
	return grpcZap.WithMessageProducer(
		func(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field) {
			span := trace.SpanFromContext(ctx)
			ctxzap.Extract(ctx).Check(level, msg).Write(
				zap.Error(err),
				zap.String("grpc.code", code.String()),
				zap.String("trace-id", span.SpanContext().TraceID().String()),
				duration,
			)
		})
}
