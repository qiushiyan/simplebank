package authgrp

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/core/user"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/business/email"
	taskcommon "github.com/qiushiyan/simplebank/business/task/common"
	"github.com/qiushiyan/simplebank/foundation/web"
	"github.com/spf13/cast"
)

type SendEmailRequest struct {
	Subject string `json:"subject" validate:"required"`
}

type SendEmailResponse struct {
	TaskId string `json:"task_id"`
}

// SendEmail godoc
//
//	@Summary		Send an email to user with given subject
//	@Description	Send an email to user with given subject, currently subject=welcome, subject=verify and subject=report are implemented. User needs to have a non-null email to be verified, and have a verified email to receive emails of other subjects.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			request	body	SendEmailRequest	true	"Email subject"
//	@Success		201
//	@Failure		400
//	@Failure		404
//	@Failure		409
//	@Router			/send-email [post]
func (h *Handler) SendEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req SendEmailRequest
	var taskId string
	var err error

	if err = web.Decode(r, &req); err != nil {
		return err
	}

	username := auth.GetUsername(ctx)
	user, err := h.user.QueryByUsername(ctx, username)
	if err != nil {
		return err
	}

	if !user.Email.Valid {
		return web.NewError(
			errors.New("user does not have an email registered"),
			http.StatusConflict,
		)
	}

	switch req.Subject {
	case email.SubjectReport:
		taskId, err = h.sendReportEmail(ctx, user)
	case email.SubjectWelcome:
		taskId, err = h.sendWelcomeEmail(ctx, user)
	case email.SubjectVerify:
		taskId, err = h.sendVerifyEmail(ctx, user)
	default:
		return web.NewError(
			email.ErrInvalidSubject,
			http.StatusNotImplemented,
		)
	}

	if err != nil {
		return err
	}

	res := SendEmailResponse{
		TaskId: taskId,
	}

	return web.RespondJson(ctx, w, res, http.StatusCreated)
}

func (h *Handler) sendReportEmail(ctx context.Context, user db_generated.User) (string, error) {
	if !user.IsEmailVerified {
		return "", web.NewError(errors.New("user's email is not verified"), http.StatusConflict)
	}
	data := email.SubjectReportData{Username: user.Username}
	payload := email.SenderPayload{
		To:      user.Email.String,
		Subject: email.SubjectReport,
		Data:    data,
	}

	return h.task.CreateTask(
		ctx,
		taskcommon.TypeEmailDelivery,
		payload,
	)
}

func (h *Handler) sendWelcomeEmail(
	ctx context.Context,
	user db_generated.User,
) (string, error) {
	if !user.IsEmailVerified {
		return "", web.NewError(errors.New("user's email is not verified"), http.StatusConflict)
	}

	data := email.SubjectWelcomeData{Username: user.Username}
	payload := email.SenderPayload{
		To:      user.Email.String,
		Subject: email.SubjectWelcome,
		Data:    data,
	}

	return h.task.CreateTask(
		ctx,
		taskcommon.TypeEmailDelivery,
		payload,
	)
}

func (h *Handler) sendVerifyEmail(
	ctx context.Context,
	u db_generated.User,
) (string, error) {
	if u.IsEmailVerified {
		return "", web.NewError(errors.New("user's email is already verified"), http.StatusConflict)
	}
	record, err := h.user.CreateVerifyEmail(ctx, u, user.NewVerifyEmail{Email: u.Email.String})
	if err != nil {
		return "", err
	}

	link, err := url.Parse(h.frontendHost + "/verify-email")
	if err != nil {
		return "", err
	}
	link.RawQuery = url.Values{
		"id":   []string{cast.ToString(record.ID)},
		"code": []string{record.SecretCode},
	}.Encode()

	data := email.SubjectVerifyData{
		Username: u.Username,
		Link:     link.String(),
	}

	payload := email.SenderPayload{
		To:      u.Email.String,
		Subject: email.SubjectVerify,
		Data:    data,
	}

	return h.task.CreateTask(ctx, taskcommon.TypeEmailDelivery, payload)
}
