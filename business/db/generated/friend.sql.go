// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: friend.sql

package db_generated

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createFriend = `-- name: CreateFriend :one
INSERT INTO friendships (from_account_id, to_account_id)
VALUES ($1, $2)
RETURNING id, from_account_id, to_account_id, status, created_at
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
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getFriend = `-- name: GetFriend :one
SELECT id, from_account_id, to_account_id, status, created_at
FROM friendships
WHERE id = $1
`

func (q *Queries) GetFriend(ctx context.Context, id int64) (Friendship, error) {
	row := q.db.QueryRow(ctx, getFriend, id)
	var i Friendship
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const listFriends = `-- name: ListFriends :many
SELECT id, from_account_id, to_account_id, status, created_at
FROM friendships
WHERE (
        $3::BIGINT IS NULL
        OR from_account_id = $3::BIGINT
    )
    AND (
        $4::BIGINT IS NULL
        OR to_account_id = $4::BIGINT
    )
    AND (
        $5::VARCHAR IS NULL
        OR status = $5::VARCHAR
    )
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListFriendsParams struct {
	Limit         int32       `json:"limit"`
	Offset        int32       `json:"offset"`
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	ToAccountID   pgtype.Int8 `json:"to_account_id"`
	Status        pgtype.Text `json:"status"`
}

func (q *Queries) ListFriends(ctx context.Context, arg ListFriendsParams) ([]Friendship, error) {
	rows, err := q.db.Query(ctx, listFriends,
		arg.Limit,
		arg.Offset,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Status,
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
			&i.Status,
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

const updateFriend = `-- name: UpdateFriend :one
UPDATE friendships
SET status = $2
WHERE id = $1
RETURNING id, from_account_id, to_account_id, status, created_at
`

type UpdateFriendParams struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func (q *Queries) UpdateFriend(ctx context.Context, arg UpdateFriendParams) (Friendship, error) {
	row := q.db.QueryRow(ctx, updateFriend, arg.ID, arg.Status)
	var i Friendship
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}
