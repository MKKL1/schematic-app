package command

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/domain/category"
)

type CreateCategoryVarsParams struct {
	PostId   int64
	Category string
	Values   []byte
}

type CreateCategoryVarsHandler decorator.CommandHandler[CreateCategoryVarsParams, any]

type createCategoryVarsHandler struct {
	repo     category.Repository
	provider category.SchemaProvider
}

func NewCreatePostCatValuesHandler(repo category.Repository, provider category.SchemaProvider) CreateCategoryVarsHandler {
	return createCategoryVarsHandler{repo, provider}
}

func (c createCategoryVarsHandler) Handle(ctx context.Context, params CreateCategoryVarsParams) (any, error) {
	categoryEntity, err := c.repo.FindCategoryByID(ctx, params.Category)
	if err != nil {
		return nil, err
	}

	schema, err := c.provider.GetValidator(categoryEntity.ValueDefinitions)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(params.Values, &result); err != nil {
		return nil, err
	}

	err = schema.ValidateData(result)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
