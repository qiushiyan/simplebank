-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING *;
-- name: ListTransfers :many
SELECT *
FROM transfers
WHERE from_account_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
-- name: GetTransfer :one
SELECT *
FROM transfers
WHERE id = $1;