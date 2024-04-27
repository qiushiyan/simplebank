package main

import (
	"context"
	"errors"
	"expvar"
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
	"github.com/qiushiyan/simplebank/foundation/web"
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
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	ctx := context.Background()
	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "bank-api", traceIDFn, events)

	err := godotenv.Load()
	if err != nil {
		log.Warn(ctx, "no .env file found")
	}

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "version", build)
	defer log.Info(ctx, "shutdown complete")

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s,mask"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			FrontendHost    string        `conf:"default:http://localhost:3001"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
		DB struct {
			URL      string `conf:"default:postgres://postgres:postgres@localhost:5432/bank,mask"`
			SslMode  string `conf:"default:disable"`
			MaxConns int    `conf:"default:4"`
		}
		Task struct {
			Manager             string `conf:"default:simple,help:\"simple\" for using native goroutines or \"asynq\" for using the [asynq](https://github.com/hibiken/asynq) library. If using asynq must also set redis url"`
			RedisUrl            string `conf:"default:localhost:6379,mask"`
			EmailSenderAddress  string `conf:"help:gmail account to send email"`
			EMAILSenderPassword string `conf:"help:app password for the gmail account,mask"`
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

	log.Info(ctx, "startup", "config", out)

	// -------------------------------------------------------------------------
	// Start Debug Service
	expvar.NewString("build").Set(cfg.Build)
	log.Info(ctx, "startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.StandardLibraryMux()); err != nil {
			log.Error(
				ctx,
				"shutdown",
				"status",
				"debug router closed",
				"host",
				cfg.Web.DebugHost,
				"msg",
				err,
			)
		}
	}()

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
		return fmt.Errorf("connecting to db: %w", err)
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

	switch taskOption {
	case task.OptionSimple:
		taskManager = simplemanager.New(simplemanager.Config{Log: log})
	case task.OptionAsynq:
		taskManager = asynqamanger.New(asynqamanger.Config{
			Log:            log,
			RedisAddr:      cfg.Task.RedisUrl,
			SenderAddr:     cfg.Task.EmailSenderAddress,
			SenderPassword: cfg.Task.EMAILSenderPassword,
		})
	}

	go func() {
		log.Info(ctx, "startup", "status", "task server started", "manager", cfg.Task.Manager)
		taskErrors <- taskManager.Start()
	}()
	defer taskManager.Close()

	// -------------------------------------------------------------------------
	// Start API service
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	apiErrors := make(chan error, 1)

	muxConfig := routes.Config{
		Shutdown:     shutdown,
		Log:          log,
		Store:        store,
		Task:         taskManager,
		FrontendHost: cfg.Web.FrontendHost,
		Build:        build,
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
		log.Info(
			ctx,
			"startup",
			"host",
			apiServer.Addr,
			"swagger",
			apiServer.Addr+"/swagger/index.html",
		)
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
