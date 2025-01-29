-- name: CreateCategory :one
INSERT INTO categories (id, name, value_definitions)
VALUES ($1, $2 ,$3)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE name = $1 LIMIT 1;

-- name: CreatePostCategory :one
INSERT INTO post_category_values (post_id, category_id, values)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPostCategory :one
SELECT * FROM post_category_values
WHERE post_id = $1 AND category_id = $2;

-- name: GetPostsByJSONValue :many
SELECT post_id, values
FROM post_category_values
WHERE category_id = $1
  AND values @? $2::jsonpath;

-- name: ListCategoriesForPost :many
SELECT c.*, pcv.values
FROM categories c
JOIN post_category_values pcv ON c.id = pcv.category_id
WHERE pcv.post_id = $1;