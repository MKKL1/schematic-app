// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
)

type CreateAuthorParams struct {
	Name     string
	UserID   int64
	Metadata []byte
}

const getAuthorByID = `-- name: GetAuthorByID :one
SELECT id, name, user_id, metadata FROM authors
WHERE id = $1
`

func (q *Queries) GetAuthorByID(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRow(ctx, getAuthorByID, id)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.Metadata,
	)
	return i, err
}

const getAuthorByName = `-- name: GetAuthorByName :one
SELECT id, name, user_id, metadata FROM authors
WHERE name = $1
`

func (q *Queries) GetAuthorByName(ctx context.Context, name string) (Author, error) {
	row := q.db.QueryRow(ctx, getAuthorByName, name)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.Metadata,
	)
	return i, err
}
