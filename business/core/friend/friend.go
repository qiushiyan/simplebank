package friend

import (
	"context"

	"github.com/qiushiyan/simplebank/business/data/limit"
	db "github.com/qiushiyan/simplebank/business/db/core"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
)

type Core struct {
	store db.Store
}

func NewCore(store db.Store) Core {
	return Core{store: store}
}

func (c *Core) NewRequest(
	ctx context.Context,
	nf NewFriend,
) (db_generated.Friendship, error) {

	friend, err := c.store.CreateFriend(ctx, db_generated.CreateFriendParams{
		FromAccountID: nf.FromAccountId,
		ToAccountID:   nf.ToAccountId,
	})
	if err != nil {
		return db_generated.Friendship{}, db.NewError(err)
	}
	return friend, nil
}

func (c *Core) ListRequests(
	ctx context.Context,
	filter QueryFilter,
	limiter limit.Limiter,
) ([]db_generated.Friendship, error) {
	friends, err := c.store.ListFriends(ctx, db_generated.ListFriendsParams{
		FromAccountID: db.NewInt8(filter.FromAccountId),
		ToAccountID:   db.NewInt8(filter.ToAccountId),
		Pending:       db.NewBool(filter.Pending),
		Accepted:      db.NewBool(filter.Accepted),
		Limit:         limiter.Limit,
		Offset:        limiter.Offset,
	})
	if err != nil {
		return []db_generated.Friendship{}, nil
	}

	return friends, nil
}

func (c *Core) AcceptRequest(ctx context.Context, id int64) (db_generated.Friendship, error) {
	friend, err := c.store.AcceptFriend(ctx, id)
	if err != nil {
		return db_generated.Friendship{}, db.NewError(err)
	}

	return friend, nil
}

func (c *Core) DeclineRequest(ctx context.Context, id int64) (db_generated.Friendship, error) {
	friend, err := c.store.DeclineFriend(ctx, id)
	if err != nil {
		return db_generated.Friendship{}, db.NewError(err)
	}

	return friend, nil
}
