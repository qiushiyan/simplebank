package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/qiushiyan/simplebank/foundation/logger"
	"github.com/qiushiyan/simplebank/foundation/web"
)

func Logger(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)
			log.Info(
				ctx,
				"request started",
				"method",
				r.Method,
				"path",
				r.URL.Path,
				"trace_id",
				v.TraceID,
			)

			err := handler(ctx, w, r)

			log.Info(
				ctx,
				"request completed",
				"method",
				r.Method,
				"path",
				r.URL.Path,
				"trace_id",
				v.TraceID,
				"status_code",
				v.StatusCode,
				"since",
				time.Since(v.Now),
			)

			return err
		}
		return h

	}

	return m
}
