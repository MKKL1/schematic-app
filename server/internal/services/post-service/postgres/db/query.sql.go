// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
)

const getPostById = `-- name: GetPostById :one
SELECT id, "desc", owner, author_known, author_unknown FROM post
WHERE id = $1
`

func (q *Queries) GetPostById(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRow(ctx, getPostById, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Desc,
		&i.Owner,
		&i.AuthorKnown,
		&i.AuthorUnknown,
	)
	return i, err
}
