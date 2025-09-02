package interceptors

import (
	"context"
	"fmt"

	"github.com/Justksenia/common/tracer"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ContextPropagationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		headers := md.Get(TraceHeader)
		if len(headers) == 0 {
			return handler(ctx, req)
		}

		traceID, err := trace.TraceIDFromHex(headers[0])
		if err != nil {
			return handler(ctx, req)
		}

		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
		})

		ctx = trace.ContextWithSpanContext(ctx, spanContext)

		ctx, span := tracer.StartSpan(ctx, fmt.Sprintf("call %s", info.FullMethod), trace.SpanKindClient)
		defer span.End()

		return handler(ctx, req)
	}
}
