package post

import (
	"context"
	"encoding/json"
)

type Entity struct {
	ID           int64
	Name         string
	Description  *string
	Owner        int64
	AuthorID     *int64
	CategoryVars []CategoryVarsEntity
	Tags         []string
}

type CategoryMetadata json.RawMessage

type CategoryVarsEntity struct {
	CategoryName string           `json:"categoryName"`
	Metadata     CategoryMetadata `json:"metadata"`
}

type Repository interface {
	FindById(ctx context.Context, id int64) (Entity, error)
	Create(ctx context.Context, model Entity) error
	GetCountForTag(ctx context.Context, tag string) (int64, error)
}
