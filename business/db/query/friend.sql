-- name: CreateFriend :one
INSERT INTO friendships (from_account_id, to_account_id)
VALUES ($1, $2)
RETURNING *;
-- name: AcceptFriend :one
UPDATE friendships
SET pending = FALSE,
    accepted = TRUE
WHERE id = $1
RETURNING *;
-- name: ListFriends :many
SELECT *
FROM friendships
WHERE sqlc.narg(from_account_id)::BIGINT IS NULL
    OR from_account_id = sqlc.narg(from_account_id)::BIGINT
    AND (
        sqlc.narg(to_account_id)::BIGINT IS NULL
        OR to_account_id = sqlc.narg(to_account_id)::BIGINT
    )
    AND (
        sqlc.narg(pending)::BOOLEAN IS NULL
        OR pending = sqlc.narg(pending)::BOOLEAN
    )
    AND (
        sqlc.narg(accepted)::BOOLEAN IS NULL
        OR accepted = sqlc.narg(accepted)::BOOLEAN
    )
ORDER BY id
LIMIT $1 OFFSET $2;
-- name: DeclineFriend :one
UPDATE friendships
SET pending = FALSE,
    accepted = FALSE
WHERE id = $1
RETURNING *;