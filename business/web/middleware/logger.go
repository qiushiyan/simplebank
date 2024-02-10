package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/qiushiyan/simplebank/foundation/web"
	"go.uber.org/zap"
)

func Logger(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)
			log.Infow("request started", "method", r.Method, "path", r.URL.Path, "trace_id", v.TraceID)

			err := handler(ctx, w, r)

			log.Infow("request completed", "method", r.Method, "path", r.URL.Path, "trace_id", v.TraceID, "status_code", v.StatusCode, "since", time.Since(v.Now))

			return err
		}
		return h

	}

	return m
}
