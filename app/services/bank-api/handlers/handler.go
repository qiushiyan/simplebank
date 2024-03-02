package handlers

import (
	"net/http"
	"os"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/accountgrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/authgrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/checkgrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/entrygrp"
	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/transfergrp"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/web/middleware"
	"github.com/qiushiyan/simplebank/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type MuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Store    db.Store
	Build    string
}

func NewMux(cfg MuxConfig) *web.App {
	mw := []web.Middleware{
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Panics(),
		middleware.Metrics(),
	}

	app := web.NewApp(cfg.Shutdown, mw...)
	accountHandler := accountgrp.New(cfg.Store)
	authHandler := authgrp.New(cfg.Store)
	transferHandler := transfergrp.New(cfg.Store)
	entryHandler := entrygrp.New(cfg.Store)
	checkHandler := checkgrp.New(cfg.Store, cfg.Build)

	// ==============================================================================
	// Account route group
	app.Handle(
		http.MethodGet,
		"/accounts/all",
		accountHandler.ListAll,
		middleware.Authenticate(),
		middleware.Authorize("ADMIN"),
	)
	app.Handle(http.MethodGet, "/accounts", accountHandler.List, middleware.Authenticate())
	app.Handle(http.MethodGet, "/accounts/:id", accountHandler.Get, middleware.Authenticate())
	app.Handle(http.MethodPost, "/accounts", accountHandler.Create, middleware.Authenticate())

	// ==============================================================================
	// Auth route group
	app.Handle(http.MethodPost, "/signup", authHandler.Signup)
	app.Handle(http.MethodPost, "/signin", authHandler.Signin)

	// ==============================================================================
	// Transfer route group
	app.Handle(http.MethodPost, "/transfer", transferHandler.Transfer, middleware.Authenticate())

	// ==============================================================================
	// Entires route group
	app.Handle(http.MethodGet, "/entries", entryHandler.List, middleware.Authenticate())

	// ==============================================================================
	// Health check route
	app.Handle(http.MethodGet, "/readiness", checkHandler.Readiness)
	app.Handle(http.MethodGet, "/liveness", checkHandler.Liveness)

	return app
}
