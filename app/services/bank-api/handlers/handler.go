package handlers

import (
	"net/http"
	"os"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers/account"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/business/web/middleware"
	"github.com/qiushiyan/simplebank/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Store    *db.Store
}

func APIMux(cfg APIMuxConfig) *web.App {
	mw := []web.Middleware{
		middleware.Logger(cfg.Log),
	}

	app := web.NewApp(cfg.Shutdown, mw...)

	app.Handle(http.MethodGet, "/accounts", account.ListAccounts)

	return app
}
