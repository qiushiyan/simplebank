package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/qiushiyan/simplebank/foundation/migrate"
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

	if err := migrate.NewFromFile("business/db/migration", cfg.DB.URL); err != nil {
		fmt.Println("migrate error: ", err)
		os.Exit(1)
	}

	fmt.Println("migrate success")
}
