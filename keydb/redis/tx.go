package redis

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/redis/go-redis/v9"
)

type pipeliner struct {
	redis.Pipeliner
}

func (p pipeliner) Close() error {
	return nil
}

func (i *Instance) Begin(_ context.Context) (*Instance, error) {
	if _, ok := i.client.(redis.Pipeliner); ok {
		return i, nil
	}

	p := pipeliner{
		Pipeliner: i.client.Pipeline(),
	}

	txClient := &Instance{
		KeyDBFactory: &KeyDBFactory{
			client:     p,
			serializer: i.serializer,
		},
		name: i.name,
		ttl:  i.ttl,
	}
	return txClient, nil
}

func (i *Instance) Commit(ctx context.Context) error {
	p, ok := i.client.(redis.Pipeliner)
	if !ok {
		return errors.New("no open transaction")
	}
	_, err := p.Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "commit transaction")
	}
	return nil
}

func (i *Instance) Rollback(_ context.Context) error {
	p, ok := i.KeyDBFactory.client.(redis.Pipeliner)
	if !ok {
		return nil
	}
	p.Discard()
	return nil
}
