package user

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/qiushiyan/simplebank/business/auth"
	"github.com/qiushiyan/simplebank/business/auth/token"
	db "github.com/qiushiyan/simplebank/business/db/core"
	. "github.com/qiushiyan/simplebank/business/db/generated"
)

type Core struct {
	store db.Store
}

func NewCore(store db.Store) Core {
	return Core{store: store}
}

func (u *Core) QueryByUsername(ctx context.Context, username string) (User, error) {
	user, err := u.store.GetUser(ctx, username)
	if err != nil {
		if db.IsNoRowsError(err) {
			return User{}, errors.New("user not found")
		}
		return User{}, db.NewError(err)
	}
	return user, nil
}

func (u *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	hash, err := auth.HashPassword(nu.Password)
	if err != nil {
		return User{}, err
	}

	user, err := u.store.CreateUser(ctx, CreateUserParams{
		Username:       nu.Username,
		HashedPassword: hash,
		Email:          db.NewText(&nu.Email),
	})

	if err != nil {
		return User{}, db.NewError(err)
	}

	return user, nil
}

func (u *Core) CreateToken(
	ctx context.Context,
	user User,
	nt NewToken,
) (token.Token, error) {
	if !auth.VerifyPassword(user.HashedPassword, nt.Password) {
		return token.Token{}, auth.NewAuthError("wrong password", http.StatusForbidden)
	}

	roles := []token.Role{token.RoleUser}
	if user.Username == token.RoleAdmin.Name() {
		roles = append(roles, token.RoleAdmin)
	}

	t, err := token.NewToken(
		user.Username,
		roles,
		24*7*30*time.Hour,
	)
	if err != nil {
		return token.Token{}, err
	}

	return t, nil
}
