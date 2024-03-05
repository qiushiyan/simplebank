package user

import (
	"context"
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

func (u *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	hash, err := auth.HashPassword(nu.Password)
	if err != nil {
		return User{}, err
	}

	user, err := u.store.CreateUser(ctx, CreateUserParams{
		Username:       nu.Username,
		HashedPassword: hash,
		Email:          nu.Email,
	})

	if err != nil {
		return User{}, db.NewError(err)
	}

	return user, nil
}

func (u *Core) CreateSession(
	ctx context.Context,
	ns NewSession,
) (token.Token, User, error) {
	user, err := u.store.GetUser(ctx, ns.Username)
	if err != nil {
		return token.Token{}, User{}, db.NewError(err)
	}

	if !auth.VerifyPassword(user.HashedPassword, ns.Password) {
		return token.Token{}, User{}, auth.NewAuthError("wrong password", http.StatusForbidden)
	}

	roles := []token.Role{token.RoleUser}
	if user.Username == token.RoleAdmin.Name() {
		roles = append(roles, token.RoleAdmin)
	}

	nt, err := token.NewToken(
		user.Username,
		roles,
		24*7*30*time.Hour,
	)
	if err != nil {
		return token.Token{}, User{}, err
	}

	return nt, user, nil
}
