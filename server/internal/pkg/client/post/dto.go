package post

import (
	"github.com/google/uuid"
	"time"
)

type FileProcessingState string

const (
	FilePending   FileProcessingState = "Pending"
	FileAvailable FileProcessingState = "Available"
)

type PostDto struct {
	ID          int64
	Name        string
	Description *string
	Owner       int64
	AuthorID    *int64
	Categories  []PostCategoryDto
	Tags        []string
	File        []PostFileDto
}

type PostFileDto struct {
	Hash      *string
	Name      string
	Downloads *int32
	FileSize  *int32
	UpdatedAt *time.Time
	State     FileProcessingState
}

type PostCategoryDto struct {
	Name     string
	Metadata map[string]interface{}
}

type CreatePostRequest struct {
	Name        string
	Description *string
	AuthorID    *int64
	Sub         uuid.UUID
	Categories  []CreatePostRequestCategory
	Tags        []string
	Files       []uuid.UUID
}

type CreatePostRequestCategory struct {
	Name     string
	Metadata map[string]interface{}
}
