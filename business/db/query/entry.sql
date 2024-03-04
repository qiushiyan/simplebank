-- name: CreateEntry :one
INSERT INTO entries (account_id, amount)
VALUES ($1, $2)
RETURNING *;
-- name: ListEntries :many
SELECT *
FROM entries
WHERE sqlc.narg(account_id)::bigint IS NULL
    OR account_id = sqlc.narg(account_id)::bigint
    AND (
        COALESCE(sqlc.narg(start_date)::timestamp, NULL) IS NULL
        OR created_at >= sqlc.narg(start_date)::timestamp
    )
    AND (
        COALESCE(sqlc.narg(end_date)::timestamp, NULL) IS NULL
        OR created_at <= sqlc.narg(end_date)::timestamp
    )
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
-- name: GetEntry :one
SELECT *
FROM entries
WHERE id = $1;