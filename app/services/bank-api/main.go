package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/qiushiyan/bank-api/foundation/logger"
	"go.uber.org/zap"
)

var build = "develop"

func main() {
	log, err := logger.New("bank-api")
	if err != nil {
		fmt.Sprintln("creating logger: %w", err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "BUILD-", build)

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s,mask"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			Host         string `conf:"default:localhost"`
			Port         string `conf:"default:5432"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:2"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
		Args conf.Args
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "A simple bank system, read docs at https://github.com/qiushiyan/simplebank",
		},
	}

	prefix := "BANK"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	defer log.Infow("Shutdown complete")

	// show the current config
	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config string: %w", err)
	}

	log.Infow("startup", "config", out)

	return nil

}
