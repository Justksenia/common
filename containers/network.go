package containers

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
)

func NewNetwork(ctx context.Context, name string) (testcontainers.Network, error) { //nolint: ireturn // it's ok
	return testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		ProviderType:   testcontainers.ProviderDocker,
		NetworkRequest: testcontainers.NetworkRequest{Name: name, CheckDuplicate: true},
	})
}
