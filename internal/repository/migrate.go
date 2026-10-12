package repository

import (
	"errors"
	"net/url"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migration opens its own connection; closing it does not close the API's pool.
func Migration(databaseURL, path string) (*migrate.Migrate, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	source := url.URL{Scheme: "file", Path: filepath.ToSlash(absolute)}
	return migrate.New(source.String(), databaseURL)
}

func Migrate(databaseURL, path string) error {
	m, err := Migration(databaseURL, path)
	if err != nil {
		return err
	}
	defer m.Close()
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
