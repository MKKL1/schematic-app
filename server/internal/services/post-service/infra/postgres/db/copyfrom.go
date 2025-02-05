// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: copyfrom.go

package db

import (
	"context"
)

// iteratorForCreatePost implements pgx.CopyFromSource.
type iteratorForCreatePost struct {
	rows                 []CreatePostParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreatePost) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreatePost) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ID,
		r.rows[0].Name,
		r.rows[0].Desc,
		r.rows[0].Owner,
		r.rows[0].AuthorKnown,
		r.rows[0].AuthorUnknown,
	}, nil
}

func (r iteratorForCreatePost) Err() error {
	return nil
}

func (q *Queries) CreatePost(ctx context.Context, arg []CreatePostParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"post"}, []string{"id", "name", "desc", "owner", "author_known", "author_unknown"}, &iteratorForCreatePost{rows: arg})
}
