package accountgrp

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/foundation/web"
)

func (h *Handler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := strings.Split(r.URL.Path, "/")[2]
	aid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	ret, err := h.store.GetAccount(ctx, int64(aid))

	if err != nil {
		return db.NewError(err)
	}

	return web.RespondJson(ctx, w, ret, http.StatusOK)
}
