package migrate

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewFromFile(sourceURL, databaseURL string) error {
	migration, err := migrate.New(fmt.Sprintf("file://%s", sourceURL), databaseURL)
	if err != nil {
		return fmt.Errorf("fail to create migration: %w", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("fail to migrate up: %w", err)
	}

	return nil
}
