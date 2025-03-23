package post

import "encoding/json"

type Post struct {
	ID          int64
	Name        string
	Description *string
	Owner       int64
	AuthorID    *int64
	Categories  PostCategories
	Tags        []string
}

type CategoryMetadata map[string]interface{}

type PostCategories json.RawMessage

type PostCategoriesStructured map[string]CategoryMetadata
