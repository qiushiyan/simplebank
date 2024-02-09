package handlers

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/foundation/web"
)

func ListAccounts(ctc context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.RespondJson(ctc, w, []string{"1", "2", "3"}, http.StatusOK)
}
