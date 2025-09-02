package redis

import (
	"context"
	"io"
	"time"

	"github.com/go-faster/errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNoData = errors.New("no data")
)

const (
	PersistentTTL = redis.KeepTTL
)

type (
	Client interface {
		redis.Cmdable
		io.Closer
	}

	Serializer interface {
		Marshal(v any) ([]byte, error)
		Unmarshal(data []byte, v any) error
	}
)

type KeyDBFactory struct {
	client           Client
	serializer       Serializer
	repetitionFactor int
}

func New(conf Config) (*KeyDBFactory, error) {
	var client Client
	cfg := toUniversalRedisConfig(conf)
	client = redis.NewUniversalClient(cfg)
	if cmd := client.Ping(context.Background()); cmd.Err() != nil {
		return nil, errors.Wrap(cmd.Err(), "init key db client")
	}
	return &KeyDBFactory{
		client:     client,
		serializer: jsoniter.ConfigCompatibleWithStandardLibrary,
	}, nil
}

func (k *KeyDBFactory) NewInstance(name string, ttl time.Duration) *Instance {
	return &Instance{
		KeyDBFactory: k,
		name:         name,
		ttl:          ttl,
	}
}

func (k *KeyDBFactory) Close() error {
	return k.client.Close()
}

func (k *KeyDBFactory) Client() Client {
	return k.client
}

type Instance struct {
	*KeyDBFactory
	name string
	ttl  time.Duration
}
