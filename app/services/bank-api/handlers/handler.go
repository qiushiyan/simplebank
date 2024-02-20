package handlers

import (
	"net/http"
	"os"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/accountgrp"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/web/middleware"
	"github.com/qiushiyan/simplebank/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Store    db.Store
}

func NewMux(cfg APIMuxConfig) *web.App {
	mw := []web.Middleware{
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Panics(),
		middleware.Metrics(),
	}

	app := web.NewApp(cfg.Shutdown, mw...)

	accountHandler := accountgrp.New(cfg.Store)

	app.Handle(http.MethodGet, "/accounts", accountHandler.List)
	app.Handle(http.MethodGet, "/accounts/:id", accountHandler.Get)

	return app
}
