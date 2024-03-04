package authgrp

import (
	"time"

	"github.com/qiushiyan/simplebank/business/core/user"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
)

type Handler struct {
	core user.Core
}

type UserResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func New(store db.Store) *Handler {
	return &Handler{
		core: user.NewCore(store),
	}
}

func NewUserResponse(user db_generated.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}
