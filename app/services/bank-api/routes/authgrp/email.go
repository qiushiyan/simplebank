package authgrp

import (
	"context"
	"net/http"

	"github.com/qiushiyan/simplebank/business/auth"
	taskcommon "github.com/qiushiyan/simplebank/business/task/common"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type SendEmailRequest struct {
	Subject string `json:"subject" validate:"required"`
}

type SendEmailResponse struct {
	TaskId string `json:"task_id"`
}

func (h *Handler) SendEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req SendEmailRequest
	if err := web.Decode(r, &req); err != nil {
		return err
	}

	username := auth.GetUsername(ctx)

	taskId, err := h.task.CreateTask(
		ctx,
		taskcommon.TypeEmailDelivery,
		taskcommon.NewEmailDeliveryPayload(
			username,
			req.Subject,
		),
	)
	if err != nil {
		return err
	}

	res := SendEmailResponse{
		TaskId: taskId,
	}

	return web.RespondJson(ctx, w, res, http.StatusCreated)
}
