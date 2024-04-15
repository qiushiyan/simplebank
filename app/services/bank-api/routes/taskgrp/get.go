package taskgrp

import (
	"context"
	"net/http"

	taskcommon "github.com/qiushiyan/simplebank/business/task/common"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type GetTaskResponse struct {
	Data taskcommon.State `json:"data"`
}

// GetTask godoc
//
//	@Summary		Inspect task state by ID
//	@Description	Inspect task state by ID
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{object}	GetTaskResponse
//	@Router			/task/{id} [get]
//	@Security		Bearer
//	@Failure		400
//	@Failure		404
func (h *Handler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")

	status, err := h.task.GetTaskState(id)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, status, http.StatusOK)
}
