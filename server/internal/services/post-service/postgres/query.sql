-- name: GetPostById :one
SELECT * FROM post
WHERE id = $1;

-- name: CreatePost :copyfrom
INSERT INTO post (id, "desc", owner, author_known, author_unknown) VALUES ($1, $2, $3, $4, $5);