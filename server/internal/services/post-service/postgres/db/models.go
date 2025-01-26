// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type GalleryImage struct {
	ImageID int64
	PostID  pgtype.Int8
	Desc    pgtype.Text
}

type Post struct {
	ID            int64
	Desc          pgtype.Text
	Owner         int64
	AuthorKnown   pgtype.Int8
	AuthorUnknown pgtype.Text
}
