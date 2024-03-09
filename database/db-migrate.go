package database

import (
	"embed"
	e "errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var migrations embed.FS

func UpgradeDatabase() error {
	fs, err := iofs.New(migrations, "migrations")
	if err != nil {
		return errors.Wrap(err, "failed to find migrations")
	}
	m, err := migrate.NewWithSourceInstance("iofs", fs, dbURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse migrations")
	}
	if err := m.Up(); err != nil && !e.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(m.Up(), "Failed to upgrade database")
	}
	return nil
}
