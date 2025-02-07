-- name: GetPostById :one
SELECT * FROM post
WHERE id = $1;

-- name: CreatePost :copyfrom
INSERT INTO post (id, name, "desc", owner, author_id) VALUES
($1, $2, $3, $4, $5);