// Package authgrp provides handlers for user registration and signin
package authgrp

import (
	"time"

	"github.com/qiushiyan/simplebank/business/core/user"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
	"github.com/qiushiyan/simplebank/business/task"
	"github.com/qiushiyan/simplebank/business/web/middleware"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type Handler struct {
	user user.Core
	task task.Manager
}

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func New(store db.Store, taskManager task.Manager) *Handler {
	return &Handler{
		user: user.NewCore(store),
		task: taskManager,
	}
}

func NewUserResponse(user db_generated.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Email:             user.Email.String,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func (h *Handler) Register(app *web.App) {
	app.POST("/signup", h.Signup)
	app.POST("/signin", h.Signin)
	app.POST("/send-email", h.SendEmail, middleware.Authenticate())
}

var _ web.RouteGroup = (*Handler)(nil)
