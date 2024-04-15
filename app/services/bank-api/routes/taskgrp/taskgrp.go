// Package taskgrp provides handlers for inspecting tasks
package taskgrp

import (
	"net/http"

	"github.com/qiushiyan/simplebank/business/task"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type Handler struct {
	task task.Manager
}

func New(task task.Manager) *Handler {
	return &Handler{task: task}
}

func (h *Handler) Register(a *web.App) {
	a.Handle(http.MethodGet, "/task/:id", h.Get)
}
