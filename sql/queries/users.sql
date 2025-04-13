-- name: CreateUser :one
INSERT INTO "user" (id, email, password_hash, last_login_at, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = $1;
