package migration

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-faster/errors"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib" // driver's import
)

const (
	defaultDialect = "pgx"
)

//nolint:gochecknoglobals // it's ok
var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory with migration files")
)

func Run(cfg *Config) error {
	flags.Usage = usage
	err := flags.Parse(os.Args[1:])
	if err != nil {
		return errors.Wrap(err, "parse flags")
	}

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return errors.New("not enough arguments")
	}

	_, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	command := args[0]

	db, err := goose.OpenDBWithDriver(defaultDialect, cfg.Postgres.DSN)
	if err != nil {
		return errors.Wrap(err, "open db")
	}
	defer db.Close()

	if err = goose.Run(command, db, *dir, args[1:]...); err != nil {
		return errors.Wrap(err, "run goose command")
	}
	return nil
}

//nolint:forbidigo //it's ok
func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

//nolint:gochecknoglobals // it's ok
var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)
