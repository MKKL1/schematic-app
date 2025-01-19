-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByOIDCSub :one
SELECT * FROM users
WHERE oidc_sub = $1;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = $1;

-- name: UserExists :one
SELECT * FROM users
WHERE name = $1;

-- name: CreateUser :copyfrom
INSERT INTO users (id, name, oidc_sub) VALUES ($1, $2, $3);