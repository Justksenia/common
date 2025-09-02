package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/adstail/ts-common/containers"
)

type RedisTestSuite struct {
	suite.Suite
	containers *TestDatabase
	instance   *Instance
}

func (s *RedisTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := SetupTestDatabase(ctx, s.T())
	if err != nil {
		s.T().Fatal(err)
	}
	s.containers = db
	s.instance = db.Factory.NewInstance("test", time.Minute)
}

func TestChannelsTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

type TestDatabase struct {
	Factory   *KeyDBFactory
	container *containers.RedisContainer
}

func SetupTestDatabase(ctx context.Context, t *testing.T) (*TestDatabase, error) {
	t.Helper()
	redis, err := containers.NewRedis(ctx, containers.RedisConf{})
	require.NoError(t, err)
	t.Cleanup(func() { _ = redis.Container.Terminate(ctx) })

	redisConf := Config{
		Addresses: []string{redis.External},
	}

	client, err := New(redisConf)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	return &TestDatabase{
		Factory:   client,
		container: redis,
	}, nil
}
