// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

type Category struct {
	ID               int64
	Name             string
	ValueDefinitions []byte
}

type PostCategoryValue struct {
	PostID     int64
	CategoryID int64
	Values     []byte
}
