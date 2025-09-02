package containers

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisConf struct {
	Image   string
	Name    string
	Network string
}

type RedisContainer struct {
	Container testcontainers.Container
	External  string
	Internal  string
}

func NewRedis(ctx context.Context, conf RedisConf) (*RedisContainer, error) {
	const (
		defaultImageName = "reg.telespace.systems:5000/base/redis:7.2.3-alpine3.18"
		defaultPort      = "6379"
	)

	containerReq := testcontainers.ContainerRequest{
		Image:        defaultImageName,
		ExposedPorts: []string{defaultPort},
		WaitingFor:   wait.ForExposedPort(),
	}

	if conf.Network != "" {
		containerReq.Networks = []string{conf.Network}
		containerReq.NetworkAliases = map[string][]string{
			conf.Network: {"redis-test"},
		}
	}

	if conf.Image != "" {
		containerReq.Image = conf.Image
	}

	if conf.Name != "" {
		containerReq.Name = conf.Name
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Logger:           testcontainers.Logger,
		Started:          true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start container")
	}

	mappedPort, err := container.MappedPort(ctx, defaultPort)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get exposed port for redis container")
	}

	networkIP, err := container.ContainerIP(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get container IP")
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get container host")
	}

	return &RedisContainer{
		Container: container,
		External:  fmt.Sprintf("%s:%s", host, mappedPort.Port()),
		Internal:  fmt.Sprintf("%s:%s", networkIP, defaultPort),
	}, nil
}
