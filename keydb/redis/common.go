package redis

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"

	"gitlab.com/adstail/ts-common/tracer"
)

func (i *Instance) Set(ctx context.Context, key string, value any) error {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	b, err := i.serializer.Marshal(value)
	if err != nil {
		return span.Error(errors.Wrap(err, "marshal"))
	}

	cmd := i.client.Set(ctx, key, b, i.ttl)
	if cmd.Err() != nil {
		return span.Error(errors.Wrap(err, "redis.Set"))
	}
	return nil
}

func (i *Instance) Get(ctx context.Context, key string, value any) error {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	cmd := i.client.Get(ctx, key)
	b, err := cmd.Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNoData
		}
		return span.Error(errors.Wrap(err, "redis.Get"))
	}

	if err = i.serializer.Unmarshal(b, value); err != nil {
		return span.Error(errors.Wrap(err, "unmarshal"))
	}
	return nil
}

func (i *Instance) IsExist(ctx context.Context, key string) (bool, error) {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	cmd := i.client.Exists(ctx, key)
	exists, err := cmd.Result()
	if err != nil {
		return false, span.Error(errors.Wrap(err, "redis.Exists"))
	}

	if exists > 0 {
		return true, nil
	}
	return false, nil
}

func (i *Instance) Delete(ctx context.Context, keys ...string) error {
	_, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindClient)
	span.AddAttribute("instance name", i.name)
	defer span.End()

	cmd := i.client.Del(ctx, keys...)
	if err := cmd.Err(); err != nil {
		return span.Error(errors.Wrap(err, "redis.Del"))
	}
	return nil
}
