-- name: CreateAuthor :copyfrom
INSERT INTO authors (
    name,
    user_id,
    metadata
) VALUES (
    $1,
    $2,
$3
);

-- name: GetAuthorByID :one
SELECT * FROM authors
WHERE id = $1;

-- name: GetAuthorByName :one
SELECT * FROM authors
WHERE name = $1;