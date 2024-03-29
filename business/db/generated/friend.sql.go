// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: friend.sql

package db_generated

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const acceptFriend = `-- name: AcceptFriend :one
UPDATE friendships
SET pending = FALSE,
    accepted = TRUE
WHERE id = $1
RETURNING id, from_account_id, to_account_id, pending, accepted, created_at
`

func (q *Queries) AcceptFriend(ctx context.Context, id int64) (Friendship, error) {
	row := q.db.QueryRow(ctx, acceptFriend, id)
	var i Friendship
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Pending,
		&i.Accepted,
		&i.CreatedAt,
	)
	return i, err
}

const createFriend = `-- name: CreateFriend :one
INSERT INTO friendships (from_account_id, to_account_id)
VALUES ($1, $2)
RETURNING id, from_account_id, to_account_id, pending, accepted, created_at
`

type CreateFriendParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
}

func (q *Queries) CreateFriend(ctx context.Context, arg CreateFriendParams) (Friendship, error) {
	row := q.db.QueryRow(ctx, createFriend, arg.FromAccountID, arg.ToAccountID)
	var i Friendship
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Pending,
		&i.Accepted,
		&i.CreatedAt,
	)
	return i, err
}

const declineFriend = `-- name: DeclineFriend :one
UPDATE friendships
SET pending = FALSE,
    accepted = FALSE
WHERE id = $1
RETURNING id, from_account_id, to_account_id, pending, accepted, created_at
`

func (q *Queries) DeclineFriend(ctx context.Context, id int64) (Friendship, error) {
	row := q.db.QueryRow(ctx, declineFriend, id)
	var i Friendship
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Pending,
		&i.Accepted,
		&i.CreatedAt,
	)
	return i, err
}

const listFriends = `-- name: ListFriends :many
SELECT id, from_account_id, to_account_id, pending, accepted, created_at
FROM friendships
WHERE $3::BIGINT IS NULL
    OR from_account_id = $3::BIGINT
    AND (
        $4::BIGINT IS NULL
        OR to_account_id = $4::BIGINT
    )
    AND (
        $5::BOOLEAN IS NULL
        OR pending = $5::BOOLEAN
    )
    AND (
        $6::BOOLEAN IS NULL
        OR accepted = $6::BOOLEAN
    )
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListFriendsParams struct {
	Limit         int32       `json:"limit"`
	Offset        int32       `json:"offset"`
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	ToAccountID   pgtype.Int8 `json:"to_account_id"`
	Pending       pgtype.Bool `json:"pending"`
	Accepted      pgtype.Bool `json:"accepted"`
}

func (q *Queries) ListFriends(ctx context.Context, arg ListFriendsParams) ([]Friendship, error) {
	rows, err := q.db.Query(ctx, listFriends,
		arg.Limit,
		arg.Offset,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Pending,
		arg.Accepted,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Friendship{}
	for rows.Next() {
		var i Friendship
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Pending,
			&i.Accepted,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
