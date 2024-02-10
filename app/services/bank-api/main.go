package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers"
	"github.com/qiushiyan/simplebank/business/web/debug"
	"github.com/qiushiyan/simplebank/foundation/logger"
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

	if err := run(context.Background(), log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(ctx context.Context, log *zap.SugaredLogger) error {

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

	// show the current config
	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config string: %w", err)
	}

	log.Infow("startup", "with config", out)

	// -------------------------------------------------------------------------
	// Start Debug Service

	log.Infow("startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.StandardLibraryMux()); err != nil {
			log.Errorw(
				"shutdown",
				"status",
				"debug router closed",
				"host",
				cfg.Web.DebugHost,
				"ERROR",
				err,
			)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// -------------------------------------------------------------------------
	// Start API service

	serverErrors := make(chan error, 1)
	apiMux := handlers.APIMux(
		handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      log,
			Store:    nil,
		},
	)
	apiServer := http.Server{
		Addr:         cfg.Web.APIHost,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		Handler:      apiMux,
	}

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", apiServer.Addr)
		serverErrors <- apiServer.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		_, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := apiServer.Shutdown(ctx); err != nil {
			apiServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil

}
