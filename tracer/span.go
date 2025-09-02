package tracer

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Это обёртка над trace.Span от OpenTelemetry.
type Span struct {
	trace.Span
}

// запускает новый span через otel.Tracer(tracerName).Start(...).
//возвращает обновлённый context.Context (чтобы потом прокидывать в другие вызовы) и обёрнутый *Span.
func StartSpan(ctx context.Context, name string, kind trace.SpanKind, opts ...trace.SpanStartOption) (context.Context, *Span) {
	opts = append(opts, trace.WithSpanKind(kind))

	ctx, span := otel.Tracer(tracerName).Start(ctx, name, opts...)
	return ctx, &Span{
		span,
	}
}

//Возвращает имя функции-вызывателя. То есть можно сделать:
func AutoFillName() string {
	return prettifier.FuncName(1)
}

//обёртка: сразу и записывает ошибку, и помечает span как Error.
func (s *Span) Error(err error) error {
	s.RecordError(err)
	s.SetStatus(codes.Error, err.Error())
	return err
}

//упрощает запись метаданных в span.
func (s *Span) AddAttribute(key string, value interface{}) {
	var attr attribute.KeyValue
	switch value := (value).(type) {
	case uintptr:
		attr = attribute.Int64(key, int64(value))
	case uint:
		attr = attribute.Int64(key, int64(value))
	case uint8:
		attr = attribute.Int64(key, int64(value))
	case uint16:
		attr = attribute.Int64(key, int64(value))
	case uint32:
		attr = attribute.Int64(key, int64(value))
	case uint64:
		attr = attribute.Int64(key, int64(value))
	case int:
		attr = attribute.Int64(key, int64(value))
	case int8:
		attr = attribute.Int64(key, int64(value))
	case int16:
		attr = attribute.Int64(key, int64(value))
	case int32:
		attr = attribute.Int64(key, int64(value))
	case int64:
		attr = attribute.Int64(key, value)
	case float32:
		attr = attribute.Float64(key, float64(value))
	case float64:
		attr = attribute.Float64(key, value)
	case string:
		attr = attribute.String(key, value)
	default:
		var str string
		data, err := json.MarshalIndent(value, "", "\t")
		if err != nil {
			str = fmt.Sprintf("ERROR AddAttrbute: %s", err.Error())
		} else {
			str = string(data)
		}
		attr = attribute.String(key, str)
	}

	s.SetAttributes(attr)
}
