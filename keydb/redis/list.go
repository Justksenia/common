package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-faster/errors"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"

	"github.com/Justksenia/common/tracer"
)

const (
	ListElementFirstPosition int64 = 0
	ListElementLastPosition  int64 = -1
)

func (i *Instance) RPush(ctx context.Context, key string, value any) error {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	b, err := i.serializer.Marshal(value)
	if err != nil {
		return span.Error(errors.Wrap(err, "marshal"))
	}
	cmd := i.client.RPush(ctx, key, string(b))
	if err = cmd.Err(); err != nil {
		return span.Error(errors.Wrap(err, "redis.RPush"))
	}

	return nil
}

func (i *Instance) LPush(ctx context.Context, key string, value any) error {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	b, err := i.serializer.Marshal(value)
	if err != nil {
		return span.Error(errors.Wrap(err, "marshal"))
	}
	cmd := i.client.LPush(ctx, key, string(b))
	if err = cmd.Err(); err != nil {
		return span.Error(errors.Wrap(err, "redis.LPush"))
	}

	return nil
}

func (i *Instance) RPop(ctx context.Context, key string, value any) error {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	cmd := i.client.RPop(ctx, key)
	b, err := cmd.Bytes()
	if err != nil {
		return span.Error(errors.Wrap(err, "redis.RPop"))
	}

	if err = i.serializer.Unmarshal(b, &value); err != nil {
		return span.Error(errors.Wrap(err, "unmarshal"))
	}
	return nil
}

func (i *Instance) LPop(ctx context.Context, key string, value any) error {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	cmd := i.client.LPop(ctx, key)
	b, err := cmd.Bytes()
	if err != nil {
		return span.Error(errors.Wrap(err, "redis.LPop"))
	}

	if err = i.serializer.Unmarshal(b, &value); err != nil {
		return span.Error(errors.Wrap(err, "unmarshal"))
	}
	return nil
}

func (i *Instance) GetList(ctx context.Context, key string, val any) error {
	const (
		startScanPosition, stopScanPosition int64 = 0, -1
	)
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	cmd := i.client.LRange(ctx, key, startScanPosition, stopScanPosition)
	res, err := cmd.Result()
	if err != nil {
		return span.Error(errors.Wrap(err, "redis.LRange"))
	}

	if len(res) == 0 {
		return ErrNoData
	}

	rc := strings.Join(res, ",")
	err = i.serializer.Unmarshal([]byte(fmt.Sprintf("[%s]", rc)), val)
	if err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	return nil
}

func (i *Instance) GetElementByPosition(ctx context.Context, key string, pos int64, val any) error {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	cmd := i.client.LIndex(ctx, key, pos)
	result, err := cmd.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNoData
		}
		return errors.Wrap(err, "redis.LIndex")
	}

	if err = i.serializer.Unmarshal([]byte(result), val); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	return nil
}

func (i *Instance) RemoveFromList(ctx context.Context, key string, value any) error {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	const (
		numberDeleteValues int64 = 1
	)

	b, err := i.serializer.Marshal(value)
	if err != nil {
		return span.Error(errors.Wrap(err, "marshal"))
	}

	cmd := i.client.LRem(ctx, key, numberDeleteValues, b)
	if err = cmd.Err(); err != nil {
		return span.Error(errors.Wrap(err, "redis.LRem"))
	}
	return nil
}
