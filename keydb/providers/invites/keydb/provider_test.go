package keydb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/Justksenia/common/containers"
	"github.com/Justksenia/common/keydb/redis"
)

type InviteLinkProviderTestSuite struct {
	suite.Suite
	instance *redis.Instance
	adapter  *InviteLinksKeyDBProvider
}

func (s *InviteLinkProviderTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := SetupTestDatabase(ctx, s.T())
	if err != nil {
		s.T().Fatal(err)
	}

	s.adapter = New(db.Factory)
	s.instance = s.adapter.client
}

func TestChannelsTestSuite(t *testing.T) {
	suite.Run(t, new(InviteLinkProviderTestSuite))
}

type TestDatabase struct {
	Factory   *redis.KeyDBFactory
	container *containers.RedisContainer
}

func SetupTestDatabase(ctx context.Context, t *testing.T) (*TestDatabase, error) {
	t.Helper()
	redisContainer, err := containers.NewRedis(ctx, containers.RedisConf{})
	require.NoError(t, err)
	t.Cleanup(func() { _ = redisContainer.Container.Terminate(ctx) })

	redisConf := redis.Config{
		Addresses: []string{redisContainer.External},
	}

	redisClient, err := redis.New(redisConf)
	require.NoError(t, err)
	t.Cleanup(func() { _ = redisClient.Close() })
	return &TestDatabase{
		Factory:   redisClient,
		container: redisContainer,
	}, nil
}
