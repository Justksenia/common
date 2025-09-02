package containers

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	setupWaitingTime       = 3 * time.Minute
	loadDatasetWaitingTIme = 10 * time.Second
)

type TestPostgresDatabase struct {
	DBInstance *bun.DB
	container  *Postgres
}

type PostgresPresetConf struct {
	PostgresConf
	DockerConf
}

func SetupPostgresTestDatabase(t *testing.T, conf PostgresPresetConf) *TestPostgresDatabase {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), setupWaitingTime)
	defer cancel()

	if conf.PostgresConf.Network == "" {
		const networkPrefix = "pg-test-net"
		conf.PostgresConf.Network = fmt.Sprintf("%s-%s", networkPrefix, t.Name())
	}

	if conf.PostgresConf.Name == "" {
		conf.PostgresConf.Name = fmt.Sprintf("pg-%s", t.Name())
	}

	network, err := NewNetwork(ctx, conf.PostgresConf.Network)
	require.NoError(t, err)
	t.Cleanup(func() { _ = network.Remove(context.Background()) })

	postgres, err := NewPostgres(ctx, conf.PostgresConf)
	require.NoError(t, err)
	t.Cleanup(func() { _ = postgres.Container.Terminate(context.Background()) })

	migration, err := NewPostgresMigrate(ctx, &MigrateConf{
		Network:     conf.PostgresConf.Network,
		PostgresURL: postgres.Internal,
		Docker:      conf.DockerConf,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = migration.Terminate(context.Background()) })

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(postgres.External + "?sslmode=disable")))
	instance := bun.NewDB(db, pgdialect.New())
	err = instance.Ping()
	require.NoError(t, err)
	t.Cleanup(func() { instance.Close() })

	return &TestPostgresDatabase{
		container:  postgres,
		DBInstance: instance,
	}
}

func (tdb *TestPostgresDatabase) LoadDataset(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), loadDatasetWaitingTIme)
	defer cancel()

	return tdb.container.LoadDataset(ctx, path)
}
