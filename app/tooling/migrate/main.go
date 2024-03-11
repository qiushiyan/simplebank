package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := struct {
		DB struct {
			URL string `conf:"default:postgres://postgres:postgres@localhost:5432/bank?sslmode=disable,mask"`
		}
		Args conf.Args
	}{}
	prefix := ""
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		}
		fmt.Println("config parse error: ", err)
		return
	}

	if err := run(cfg.DB.URL); err != nil {
		fmt.Println("migrate error: ", err)
		os.Exit(1)
	}

	fmt.Println("migrate success")

}

func run(url string) error {
	migration, err := migrate.New("file://business/db/migration", url)
	if err != nil {
		return fmt.Errorf("fail to create migration: %w", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("fail to migrate up: %w", err)
	}

	return nil
}
