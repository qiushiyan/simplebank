package routes

import (
	"os"

	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/accountgrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/authgrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/checkgrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/docsgrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/entrygrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/friendgrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/taskgrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes/transfergrp"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/task"
	"github.com/qiushiyan/simplebank/business/web/middleware"
	"github.com/qiushiyan/simplebank/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type Config struct {
	Shutdown     chan os.Signal
	Log          *zap.SugaredLogger
	Store        db.Store
	Task         task.Manager
	FrontendHost string
	Build        string
}

func NewMux(cfg Config) *web.App {
	globalMiddlewares := []web.Middleware{
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Panics(),
		middleware.Metrics(),
	}

	app := web.NewApp(cfg.Shutdown, globalMiddlewares...)

	app.AddGroup(
		accountgrp.New(cfg.Store),
		authgrp.New(cfg.Store, cfg.Task, cfg.FrontendHost),
		transfergrp.New(cfg.Store),
		entrygrp.New(cfg.Store),
		friendgrp.New(cfg.Store),
		checkgrp.New(cfg.Store, cfg.Build),
		taskgrp.New(cfg.Task),
		docsgrp.New(),
	)

	return app
}
