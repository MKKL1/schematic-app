-- name: CreateImage :exec
INSERT INTO image (file_hash, image_type)
VALUES ($1, $2);

-- name: GetImageTypesForHash :many
SELECT image_type FROM image WHERE file_hash = $1;