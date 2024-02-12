package middleware

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/web/response"
	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/qiushiyan/simplebank/foundation/web"
	"go.uber.org/zap"
)

func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {

				log.Errorw("ERROR", "trace_id", web.GetTraceID(ctx), "message", err)

				var er response.ErrorDocument
				var status int
				// trusted error
				switch {
				case response.IsError(err):
					reqErr := response.GetError(err)
					if validate.IsFieldErrors(reqErr.Err) {
						fieldErrors := validate.GetFieldErrors(reqErr.Err)
						er = response.ErrorDocument{
							Error:  "data validation error",
							Fields: fieldErrors.Fields(),
						}
						status = reqErr.Status
						break
					}

					er = response.ErrorDocument{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				default:
					er = response.ErrorDocument{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.RespondJson(ctx, w, er, status); err != nil {
					return err
				}

				// shutdown error is also a kinda trusted error, but we'll leave it to the root handle method
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
