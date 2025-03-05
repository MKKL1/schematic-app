-- name: CreateFile :one
INSERT INTO tmp_file (file_hash, file_name, content_type, file_size, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFileByHash :one
SELECT * FROM tmp_file
WHERE file_hash = $1;

-- name: ListExpiredFiles :many
SELECT file_hash, expires_at
FROM tmp_file
WHERE expires_at < NOW();