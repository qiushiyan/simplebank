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
	"github.com/qiushiyan/simplebank/foundation/random"
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

func (u *Core) CreateTx(
	ctx context.Context,
	nu NewUser,
	afterCreate db.AfterCreateUserFunc,
) (db.CreateUserTxResult, error) {

	hash, err := auth.HashPassword(nu.Password)
	if err != nil {
		return db.CreateUserTxResult{}, err
	}

	result, err := u.store.CreateUserTx(ctx, db.CreateUserTxParams{
		CreateUserParams: CreateUserParams{
			Username:       nu.Username,
			HashedPassword: hash,
			Email:          db.NewText(&nu.Email),
		},
		AfterCreate: afterCreate,
	})

	if err != nil {
		return db.CreateUserTxResult{}, db.NewError(err)
	}

	return result, nil
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

// CreateVerifyEmail creates a new verify email record for the user, returns the record id, secret code and error
func (u *Core) CreateVerifyEmail(
	ctx context.Context,
	user User,
	ne NewVerifyEmail,
) (VerifyEmail, error) {
	code := random.RandomString(6)
	record, err := u.store.CreateVerifyEmail(ctx, CreateVerifyEmailParams{
		Username:   user.Username,
		SecretCode: code,
		Email:      ne.Email,
	})

	if err != nil {
		return VerifyEmail{}, db.NewError(err)
	}

	return record, nil
}

func (u *Core) VerifyEmail(
	ctx context.Context,
	user User,
	fe FinishVerifyEmail,
) (bool, error) {
	err := u.store.ExecuteInTransaction(ctx, func(q *Queries) error {
		_, err := u.store.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         fe.Id,
			SecretCode: fe.Code,
		})

		if err != nil {
			if db.IsNoRowsError(err) {
				return errors.New("no verification record exists or it has expired")
			}

			return err
		}

		var val = true
		_, err = u.store.UpdateUser(ctx, UpdateUserParams{
			Username:        user.Username,
			IsEmailVerified: db.NewBool(&val),
		})

		return err
	})
	if err != nil {
		return false, db.NewError(err)
	}

	return true, nil
}
