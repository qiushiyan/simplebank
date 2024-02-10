package middleware

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/foundation/web"
	"go.uber.org/zap"
)

func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Errorw("ERROR", "trace_id", web.GetTraceID(ctx), "message", err)

				switch {

				}

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
