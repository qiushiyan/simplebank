package middleware

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/foundation/web"
	"go.uber.org/zap"
)

func Logger(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			log.Infow("request started", "method", r.Method, "path", r.URL.Path)

			err := handler(ctx, w, r)
			if err != nil {
				log.Errorw("error handling request", "error", err)
			}

			log.Infow("request completed", "method", r.Method, "path", r.URL.Path)

			return err
		}
		return h

	}

	return m
}
