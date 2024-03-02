package transfergrp

import db "github.com/qiushiyan/simplebank/business/db/core"

type Handler struct {
	store db.Store
}

func New(store db.Store) *Handler {
	return &Handler{store: store}
}
