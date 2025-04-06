-- name: CreateTempFile :exec
INSERT INTO tmp_file (store_key, file_name, content_type, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetTempFile :one
SELECT * FROM tmp_file
WHERE store_key = $1;

-- name: ListExpiredFiles :many
SELECT store_key, expires_at
FROM tmp_file
WHERE expires_at < NOW();

-- name: DeleteTmpFiles :exec
DELETE FROM tmp_file
WHERE store_key = ANY($1::text[]);

-- name: FileExistsByHash :one
SELECT exists(SELECT 1 FROM file
              WHERE hash = $1);

-- name: CreateFile :exec
INSERT INTO file (hash, file_size, content_type)
VALUES ($1, $2, $3);

-- name: GetAndMarkTempFileProcessing :one
UPDATE tmp_file
SET status = 'processing', updated_at = NOW()
WHERE store_key = $1 AND status = 'pending'
RETURNING *;

-- name: MarkTempFileFailed :exec
UPDATE tmp_file
SET status = 'failed',
    processing_attempts = processing_attempts + 1,
    updated_at = NOW()
WHERE store_key = $1;

-- name: MarkTempFileProcessed :exec
UPDATE tmp_file
SET status = 'processed',
    final_hash = $2,
    updated_at = NOW()
WHERE store_key = $1;