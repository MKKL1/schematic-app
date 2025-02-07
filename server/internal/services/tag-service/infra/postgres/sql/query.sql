-- name: CreateCategory :one
INSERT INTO categories (name, value_definitions)
VALUES ($1, $2)
RETURNING *;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE name = $1 LIMIT 1;

-- name: CreatePostCategory :exec
INSERT INTO post_category_values (post_id, category, values)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPostCategory :one
SELECT * FROM post_category_values
WHERE post_id = $1 AND category = $2;

-- name: GetPostsByJSONValue :many
SELECT post_id, values
FROM post_category_values
WHERE category = $1
  AND values @? $2::jsonpath;

-- name: GetCategVarsForPost :many
SELECT category, values FROM post_category_values
WHERE post_id = $1;