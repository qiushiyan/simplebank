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
	var s *string
	if filter.Status != nil {
		s = &filter.Status.name
	}
	friends, err := c.store.ListFriends(ctx, db_generated.ListFriendsParams{
		FromAccountID: db.NewInt8(filter.FromAccountId),
		ToAccountID:   db.NewInt8(filter.ToAccountId),
		Status:        db.NewText(s),
		Limit:         limiter.Limit,
		Offset:        limiter.Offset,
	})
	if err != nil {
		return []db_generated.Friendship{}, nil
	}

	return friends, nil
}

func (c *Core) GetFriendRequest(
	ctx context.Context,
	id int64,
) (db_generated.Friendship, error) {
	friend, err := c.store.GetFriend(ctx, id)
	if err != nil {
		return db_generated.Friendship{}, db.NewError(err)
	}

	return friend, nil
}

func (c *Core) UpdateFriendRequest(
	ctx context.Context,
	id int64,
	status Status,
) (db_generated.Friendship, error) {
	friend, err := c.store.UpdateFriend(ctx, db_generated.UpdateFriendParams{
		ID:     id,
		Status: status.name,
	})
	if err != nil {
		return db_generated.Friendship{}, db.NewError(err)
	}

	return friend, nil
}
