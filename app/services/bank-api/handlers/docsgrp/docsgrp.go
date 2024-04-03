// Package docsgrp serves the swagger documentation at /swagger/index.html
package docsgrp

import (
	"context"
	"net/http"

	_ "github.com/qiushiyan/simplebank/app/services/bank-api/docs"
	"github.com/qiushiyan/simplebank/foundation/web"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (handler *Handler) Register(app *web.App) {
	h := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)
	app.GET("/swagger/*", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTP(w, r)
		return nil
	})
	app.GET("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return web.RespondJsonPlain(ctx, w, struct {
			Message string `json:"message"`
		}{
			Message: "Welcome to the Simple Bank API, see swagger docs at /swagger/index.html and development instructions at https://github.com/qiushiyan/simplebank.",
		}, http.StatusOK)
	})
}
