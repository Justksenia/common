//nolint:nosprintfhostport //it's ok
package containers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-faster/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresConf struct {
	ImageAddress string
	Name         string
	Network      string
}

type Postgres struct {
	Container testcontainers.Container
	External  string
	Internal  string
}

func NewPostgres(ctx context.Context, conf PostgresConf) (*Postgres, error) {
	const (
		defaultImageAddress = "reg.telespace.systems:5000/base/postgres:14-alpine"
		defaultPort         = "5432/tcp"
	)
	const (
		pgUser     = "admin"
		pgPassword = "password"
		pgDB       = "stats"
	)
	var env = map[string]string{
		"POSTGRES_USER":     pgUser,
		"POSTGRES_PASSWORD": pgPassword,
		"POSTGRES_DB":       pgDB,
	}
	containerReq := testcontainers.ContainerRequest{
		Image:        defaultImageAddress,
		Env:          env,
		ExposedPorts: []string{defaultPort},
		WaitingFor:   wait.ForExposedPort(),
	}

	if conf.Network != "" {
		containerReq.Networks = []string{conf.Network}
		containerReq.NetworkAliases = map[string][]string{
			conf.Network: {"postgres-test"},
		}
	}

	if conf.ImageAddress != "" {
		containerReq.Image = conf.ImageAddress
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
		return nil, errors.Wrap(err, "start container")
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, errors.Wrap(err, "get exposed port for container port 5432")
	}

	networkIP, err := container.ContainerIP(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get container IP")
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get container host")
	}

	return &Postgres{
		Container: container,
		External:  fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPassword, host, mappedPort.Port(), pgDB),
		Internal:  fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPassword, networkIP, "5432", pgDB),
	}, nil
}

func (p *Postgres) LoadDataset(ctx context.Context, path string) error {
	const timeout = 10 * time.Second

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(p.External + "?sslmode=disable")))
	conn := bun.NewDB(db, pgdialect.New())
	if err := conn.Ping(); err != nil {
		return errors.Wrap(err, "connect to instance")
	}

	c, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if _, err = conn.ExecContext(ctx, string(c)); err != nil {
		return err
	}
	return nil
}
