// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

type Category struct {
	Name             string
	ValueDefinitions []byte
}

type PostCategoryValue struct {
	PostID   int64
	Category string
	Values   []byte
}
