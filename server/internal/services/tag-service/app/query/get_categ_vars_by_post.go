package query

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/domain/category"
)

type GetCategVarsByPostHandler decorator.QueryHandler[int64, []category.PostCategoryVars]

type getCategVarsByPostHandler struct {
	repo category.Repository
}

func NewGetCategVarsByPost(repo category.Repository) GetCategVarsByPostHandler {
	return getCategVarsByPostHandler{repo}
}

// TODO entity is returned here, bad idea
func (c getCategVarsByPostHandler) Handle(ctx context.Context, id int64) ([]category.PostCategoryVars, error) {
	vars, err := c.repo.FindCategVarsByPostID(ctx, id)
	if err != nil {
		return nil, err
	}

	return vars, nil
}
