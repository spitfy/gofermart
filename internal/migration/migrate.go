package migration

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spitfy/gofermart/internal/config"
	"path/filepath"
)

type Migration struct {
	conf *config.Config
}

func NewMigration(conf *config.Config) *Migration {
	return &Migration{conf}
}

func (mg *Migration) Up() error {
	migrationsPath, err := getMigrationsDir()
	if err != nil {
		return fmt.Errorf("%w; migrations path source error", err)
	}
	migrationSource := "file://" + filepath.ToSlash(migrationsPath)
	m, err := migrate.New(
		migrationSource,
		mg.conf.DB.DatabaseURI)
	if err != nil {
		return err
	}
	defer func() {
		serr, dberr := m.Close()
		if serr != nil {
			if err != nil {
				err = fmt.Errorf("%w; migration source error: %v", err, serr)
			} else {
				err = serr
			}
		}
		if dberr != nil {
			if err != nil {
				err = fmt.Errorf("%w; migration database error: %v", err, dberr)
			} else {
				err = dberr
			}
		}
	}()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%w; migration up error", err)
	}
	return nil
}
