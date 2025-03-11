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