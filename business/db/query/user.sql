-- name: CreateUser :one
INSERT INTO users (username, hashed_password)
VALUES ($1, $2)
RETURNING *;
-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;
-- name: UpdateUser :one
UPDATE users
SET hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    password_changed_at = COALESCE(
        sqlc.narg(password_changed_at),
        password_changed_at
    )
WHERE username = sqlc.arg(username)
RETURNING *;