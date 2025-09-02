package containers

import (
	"context"
	"log"
	"os"

	"github.com/go-faster/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DockerConf struct {
	Context  string
	FileName string
}

type MigrateConf struct {
	Network     string
	PostgresURL string
	Docker      DockerConf
}

type PostgresMigrate struct {
	testcontainers.Container
}

func NewPostgresMigrate(ctx context.Context, deps *MigrateConf) (*PostgresMigrate, error) {
	nameToken := os.Getenv("GIT_NAME_TOKEN")
	if nameToken == "" {
		return nil, errors.New("GIT_NAME_TOKEN is empty")
	}

	gitToken := os.Getenv("GIT_TOKEN")
	if gitToken == "" {
		return nil, errors.New("GIT_TOKEN is empty")
	}

	if deps.Docker.Context == "" {
		deps.Docker.Context = "../../../../../"
	}

	if deps.Docker.FileName == "" {
		deps.Docker.FileName = "Dockerfile.migrate"
	}

	containerReq := testcontainers.ContainerRequest{
		Cmd:            []string{"up"},
		NetworkAliases: map[string][]string{deps.Network: {"migrate"}},
		Env:            map[string]string{"POSTGRES_DSN": deps.PostgresURL},
		Networks:       []string{deps.Network},
		WaitingFor:     wait.ForExit(),
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    deps.Docker.Context,
			Dockerfile: deps.Docker.FileName,
			BuildArgs:  map[string]*string{"GIT_NAME_TOKEN": &nameToken, "GIT_TOKEN": &gitToken},
		},
	}

	if err := containerReq.Validate(); err != nil {
		return nil, errors.Wrap(err, "validate migrate container request")
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Logger:           log.New(os.Stderr, "", log.LstdFlags),
		Started:          true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "start migrate container")
	}

	return &PostgresMigrate{Container: container}, nil
}
