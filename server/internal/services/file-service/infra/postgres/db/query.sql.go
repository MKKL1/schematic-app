// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createFile = `-- name: CreateFile :exec
INSERT INTO tmp_file (store_key, file_name, expires_at)
VALUES ($1, $2, $3)
`

type CreateFileParams struct {
	StoreKey  string
	FileName  string
	ExpiresAt pgtype.Timestamptz
}

func (q *Queries) CreateFile(ctx context.Context, arg CreateFileParams) error {
	_, err := q.db.Exec(ctx, createFile, arg.StoreKey, arg.FileName, arg.ExpiresAt)
	return err
}

const deleteExpiredFiles = `-- name: DeleteExpiredFiles :exec
DELETE FROM tmp_file
WHERE store_key = ANY($1::text[])
`

func (q *Queries) DeleteExpiredFiles(ctx context.Context, dollar_1 []string) error {
	_, err := q.db.Exec(ctx, deleteExpiredFiles, dollar_1)
	return err
}

const getFile = `-- name: GetFile :one
SELECT store_key, file_name, expires_at, created_at, updated_at FROM tmp_file
WHERE store_key = $1
`

func (q *Queries) GetFile(ctx context.Context, storeKey string) (TmpFile, error) {
	row := q.db.QueryRow(ctx, getFile, storeKey)
	var i TmpFile
	err := row.Scan(
		&i.StoreKey,
		&i.FileName,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listExpiredFiles = `-- name: ListExpiredFiles :many
SELECT store_key, expires_at
FROM tmp_file
WHERE expires_at < NOW()
`

type ListExpiredFilesRow struct {
	StoreKey  string
	ExpiresAt pgtype.Timestamptz
}

func (q *Queries) ListExpiredFiles(ctx context.Context) ([]ListExpiredFilesRow, error) {
	rows, err := q.db.Query(ctx, listExpiredFiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListExpiredFilesRow
	for rows.Next() {
		var i ListExpiredFilesRow
		if err := rows.Scan(&i.StoreKey, &i.ExpiresAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
