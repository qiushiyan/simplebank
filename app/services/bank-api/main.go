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
	"github.com/joho/godotenv"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/task"
	asynqamanger "github.com/qiushiyan/simplebank/business/task/asynq"
	simplemanager "github.com/qiushiyan/simplebank/business/task/simple"
	"github.com/qiushiyan/simplebank/business/web/debug"
	"github.com/qiushiyan/simplebank/foundation/logger"
	"go.uber.org/zap"
)

var build = "develop"

//	@title			SimpleBank API
//	@version		1.0
//	@description	Example API for a banking system, see development notes at https://github.com/qiushiyan/simplebank

//	@contact.name	Qiushi Yan
//	@contact.url	github.com/qiushiyan/simplebank/issues
//	@contact.email	qiushi.yann@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host						localhost:3000
// @BasePath					/
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	log, err := logger.New("bank-api")
	if err != nil {
		fmt.Sprintln("creating logger: %w", err)
		os.Exit(1)
	}
	defer log.Sync()

	err = godotenv.Load()
	if err != nil {
		log.Warn("no .env file found")
	}

	if err := run(context.Background(), log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(ctx context.Context, log *zap.SugaredLogger) error {
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
			URL      string `conf:"default:postgres://postgres:postgres@localhost:5432/bank,mask"`
			SslMode  string `conf:"default:disable"`
			MaxConns int    `conf:"default:4"`
		}
		Task struct {
			Manager  string `conf:"default:simple,help:\"simple\" for using native goroutines or \"asynq\" for using the [asynq](https://github.com/hibiken/asynq) library. If using asynq must also set redis url"`
			RedisUrl string `conf:"default:localhost:6379,mask"`
		}
		Args conf.Args
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "A simple bank system, read docs at https://github.com/qiushiyan/simplebank",
		},
	}

	prefix := ""
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
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "BUILD-", build)
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
	// Connect database
	dbConfigString := fmt.Sprintf(
		"%s?sslmode=%s&pool_max_conns=%d",
		cfg.DB.URL,
		cfg.DB.SslMode,
		cfg.DB.MaxConns,
	)
	pool, err := db.NewPgxPool(ctx, dbConfigString)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewPostgresStore(pool)
	defer pool.Close()

	// -------------------------------------------------------------------------
	// Start Task Service
	taskErrors := make(chan error, 1)
	var taskManager task.Manager

	taskOption, err := task.ParseOption(cfg.Task.Manager)
	if err != nil {
		return err
	}

	if taskOption == task.OptionAsynq {
		taskManager = asynqamanger.New(log, cfg.Task.RedisUrl)
	} else if taskOption == task.OptionSimple {
		taskManager = simplemanager.New(log)
	}

	go func() {
		log.Infow("startup", "status", "task server started", "manager", cfg.Task.Manager)
		taskErrors <- taskManager.Start()
	}()
	defer taskManager.Close()

	// -------------------------------------------------------------------------
	// Start API service

	apiErrors := make(chan error, 1)
	muxConfig := routes.MuxConfig{
		Shutdown: shutdown,
		Log:      log,
		Store:    store,
		Task:     taskManager,
		Build:    build,
	}
	apiMux := routes.NewMux(muxConfig)

	apiServer := http.Server{
		Addr:         cfg.Web.APIHost,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		Handler:      apiMux,
	}

	go func() {
		log.Infow("startup", "host", apiServer.Addr)
		apiErrors <- apiServer.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-apiErrors:
		return fmt.Errorf("api server error: %w", err)
	case err := <-taskErrors:
		return fmt.Errorf("task manager error: %w", err)
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
