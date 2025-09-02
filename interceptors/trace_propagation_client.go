package interceptors

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const TraceHeader = "X-Trace-ID"

func TracePropagationUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		span := trace.SpanFromContext(ctx)

		ctx = metadata.AppendToOutgoingContext(ctx, TraceHeader, span.SpanContext().TraceID().String())
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
