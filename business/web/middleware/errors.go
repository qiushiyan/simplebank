package middleware

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
	"go.uber.org/zap"
)

func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {

				log.Errorw("ERROR", "trace_id", web.GetTraceID(ctx), "message", err)

				var er web.ErrorDocument
				var status int
				// trusted error
				switch {
				case validate.IsFieldErrors(err):
					fieldErrors := validate.GetFieldErrors(err)
					er = web.ErrorDocument{
						Error:  "malformed request data",
						Fields: fieldErrors.Fields(),
					}
					status = http.StatusBadRequest

				case auth.IsAuthError(err):
					authErr := auth.GetAuthError(err)
					er = web.ErrorDocument{
						Error: authErr.Error(),
					}
					status = authErr.Status

				case db.IsError(err):
					dbErr := db.GetError(err)
					er = web.ErrorDocument{
						Error: dbErr.Error(),
					}
					status = dbErr.Status

				case web.IsError(err):
					reqErr := web.GetError(err)
					er = web.ErrorDocument{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				default:
					er = web.ErrorDocument{
						Error: err.Error(),
					}
					status = http.StatusInternalServerError
				}

				if err := web.RespondJson(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}

			}

			return nil
		}
		return h
	}

	return m
}
