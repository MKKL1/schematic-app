-- name: CreateFile :exec
INSERT INTO tmp_file (store_key, file_name, expires_at)
VALUES ($1, $2, $3);

-- name: GetFile :one
SELECT * FROM tmp_file
WHERE store_key = $1;

-- name: ListExpiredFiles :many
SELECT store_key, expires_at
FROM tmp_file
WHERE expires_at < NOW();

-- name: DeleteExpiredFiles :exec
DELETE FROM tmp_file
WHERE store_key = ANY($1::text[]);