package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type CreatePostParams struct {
	Name        string
	Description *string
	AuthorID    *int64
	Sub         uuid.UUID
	Categories  []CategoryMetadataParams
	Tags        []string
}

type CategoryMetadataParams struct {
	Name     string
	Metadata map[string]interface{}
}

type CreatePostHandler decorator.CommandHandler[CreatePostParams, int64]

type createPostHandler struct {
	repo           post.Repository
	categoryRepo   category.Repository
	idNode         *snowflake.Node
	userService    client.UserApplication
	schemaProvider category.SchemaProvider
}

func NewCreatePostHandler(repo post.Repository, categoryRepo category.Repository, idNode *snowflake.Node, userService client.UserApplication, provider category.SchemaProvider) CreatePostHandler {
	return createPostHandler{repo, categoryRepo, idNode, userService, provider}
}

func (h createPostHandler) Handle(ctx context.Context, params CreatePostParams) (int64, error) {
	err := h.validateCategories(ctx, params.Categories)
	if err != nil {
		return 0, err
	}

	user, err := h.userService.Query.GetUserBySub(ctx, params.Sub)
	if err != nil {
		return 0, err
	}

	categs := make([]post.CreateCategMetadataParams, len(params.Categories))
	for i, c := range params.Categories {
		categs[i] = post.CreateCategMetadataParams{
			Name:     c.Name,
			Metadata: c.Metadata,
		}
	}

	newPost := post.CreatePostParams{
		ID:          h.idNode.Generate().Int64(),
		Name:        params.Name,
		Description: params.Description,
		Owner:       user.ID,
		AuthorID:    params.AuthorID,
		Categories:  categs,
		Tags:        params.Tags,
	}

	err = h.repo.Create(ctx, newPost)
	if err != nil {
		return 0, apperr.WrapErrorf(err, apperr.ErrorCodeUnknown, "repo.Create")
	}

	return newPost.ID, nil
}

func (h createPostHandler) validateCategories(ctx context.Context, categories []CategoryMetadataParams) error {
	for _, categoryData := range categories {
		//TODO can be bulk
		categ, err := h.categoryRepo.FindCategoryByName(ctx, categoryData.Name)
		if err != nil {
			return err
		}

		validator, err := h.schemaProvider.GetValidator(categ.ValueDefinitions)
		if err != nil {
			return err
		}

		err = validator.ValidateData(categoryData.Metadata)
		if err != nil {
			return err
		}
	}

	return nil
}
