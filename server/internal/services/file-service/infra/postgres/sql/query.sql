-- name: CreateFile :one
INSERT INTO tmp_file (file_hash, store_key, file_name, content_type, file_size, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (file_hash) DO UPDATE
    SET expires_at = EXCLUDED.expires_at,
        updated_at = NOW()
RETURNING store_key;

-- name: GetFileByHash :one
SELECT * FROM tmp_file
WHERE file_hash = $1;

-- name: GetFileByKey :one
SELECT * FROM tmp_file
WHERE store_key = $1;

-- name: ListExpiredFiles :many
SELECT file_hash, store_key, expires_at
FROM tmp_file
WHERE expires_at < NOW();

-- name: DeleteExpiredFilesByKey :exec
DELETE FROM tmp_file
WHERE store_key = ANY($1::text[]);