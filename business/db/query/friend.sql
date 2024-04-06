-- name: CreateFriend :one
INSERT INTO friendships (from_account_id, to_account_id)
VALUES ($1, $2)
RETURNING *;
-- name: GetFriend :one
SELECT *
FROM friendships
WHERE id = $1;
-- name: ListFriends :many
SELECT *
FROM friendships
WHERE (
        sqlc.narg(from_account_id)::BIGINT IS NULL
        OR from_account_id = sqlc.narg(from_account_id)::BIGINT
    )
    AND (
        sqlc.narg(to_account_id)::BIGINT IS NULL
        OR to_account_id = sqlc.narg(to_account_id)::BIGINT
    )
    AND (
        sqlc.narg(status)::VARCHAR IS NULL
        OR status = sqlc.narg(status)::VARCHAR
    )
ORDER BY id
LIMIT $1 OFFSET $2;
-- name: UpdateFriend :one
UPDATE friendships
SET status = $2
WHERE id = $1
RETURNING *;