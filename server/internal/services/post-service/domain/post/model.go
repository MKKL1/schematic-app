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
	Categories  PostCategories
	Tags        []string
	Files       []PostFile
}

type FileProcessingState string

const (
	FilePending   FileProcessingState = "Pending"
	FileAvailable FileProcessingState = "Available"
	FileFailed    FileProcessingState = "Failed"
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

type PostCategories json.RawMessage

type PostCategoriesStructured map[string]CategoryMetadata
