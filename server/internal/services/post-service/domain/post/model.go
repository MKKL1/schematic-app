package post

import (
	"encoding/json"
	"time"
)

type Post struct {
	ID          int64
	Name        string
	Description *string
	Owner       int64
	AuthorID    *int64
	Categories  PostCategoriesRaw
	Tags        []string
	Files       []PostFile
}

type FileProcessingState string

const (
	FilePending   FileProcessingState = "Pending"
	FileAvailable FileProcessingState = "Available"
)

type PostFile struct {
	Hash      *string
	Name      string
	Downloads *int32
	FileSize  *int32
	UpdatedAt *time.Time
	State     FileProcessingState
}

type CategoryMetadata map[string]interface{}

// PostCategoriesRaw is used when category data is passed without any modification to it
type PostCategoriesRaw json.RawMessage

// PostCategories contrary to PostCategoriesRaw it is used when category data needs to be structured
type PostCategories map[string]CategoryMetadata
