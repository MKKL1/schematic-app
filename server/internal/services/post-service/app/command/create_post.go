package command

import (
	"context"
	"errors"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category/validator"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
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
	Files       []CreatePostFileParams
}

type CategoryMetadataParams struct {
	Name     string
	Metadata map[string]interface{}
}

type CreatePostFileParams struct {
	TempId uuid.UUID
}

type CreatePostHandler decorator.CommandHandler[CreatePostParams, int64]

type createPostHandler struct {
	repo         post.Repository
	categoryRepo category.Repository
	idNode       *snowflake.Node
	userService  client.UserApplication
	eventBus     *cqrs.EventBus
}

func NewCreatePostHandler(repo post.Repository, categoryRepo category.Repository, idNode *snowflake.Node, userService client.UserApplication, eventBus *cqrs.EventBus) CreatePostHandler {
	return createPostHandler{repo, categoryRepo, idNode, userService, eventBus}
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

	categs := make([]post.CreatePostCategoryParams, len(params.Categories))
	for i, c := range params.Categories {
		categs[i] = post.CreatePostCategoryParams{
			Name:     c.Name,
			Metadata: c.Metadata,
		}
	}

	files := make([]post.CreatePostFileParams, len(params.Files))
	for i, f := range params.Files {
		files[i] = post.CreatePostFileParams{
			TempId: f.TempId,
		}
	}

	newPost := post.CreatePostParams{
		ID:          h.idNode.Generate().Int64(),
		Name:        params.Name,
		Description: params.Description,
		AuthorID:    params.AuthorID,
		Owner:       user.ID,
		Categories:  categs,
		Tags:        params.Tags,
		Files:       files,
	}

	err = h.repo.Create(ctx, newPost)
	if err != nil {
		return 0, apperr.WrapErrorf(err, apperr.ErrorCodeUnknown, "CreatePostHandler: Handle: repo.Create")
	}

	go func() {
		err = h.publishCreatePostEvent(ctx, newPost, params.Files)
		if err != nil {
			//TODO handle
			//log or delete created post from database
		}
	}()

	return newPost.ID, nil
}

func (h createPostHandler) publishCreatePostEvent(ctx context.Context, newPost post.CreatePostParams, files []CreatePostFileParams) error {
	categs := make(post.PostCategoriesStructured, len(newPost.Categories))
	for _, c := range newPost.Categories {
		categs[c.Name] = c.Metadata
	}

	eventFiles := make([]post.PostCreatedFileData, len(files))
	for i, f := range files {
		eventFiles[i] = post.PostCreatedFileData{
			TempId: f.TempId.String(),
		}
	}

	event := post.PostCreated{
		Id:          newPost.ID,
		Name:        newPost.Name,
		Description: newPost.Description,
		Owner:       newPost.Owner,
		AuthorId:    newPost.AuthorID,
		Categories:  categs,
		Tags:        newPost.Tags,
		Files:       eventFiles,
	}

	err := h.eventBus.Publish(ctx, event)
	if err != nil {
		return apperr.WrapErrorf(err, apperr.ErrorCodeUnknown, "CreatePostHandler: Publish")
	}
	return nil
}

func (h createPostHandler) validateCategories(ctx context.Context, categories []CategoryMetadataParams) error {
	verrs := make(map[string]validator.ValidationError)
	var ok = true
	for _, categoryData := range categories {
		//TODO can be bulk
		categ, err := h.categoryRepo.FindCategoryByName(ctx, categoryData.Name)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return apperr.WrapErrorf(err, category.ErrorCodeCategoryNotFound, "CreatePostHandler: validateCategories: repo.FindCategoryByName").
					AddMetadata("name", categoryData.Name)
			}
			return apperr.WrapErrorf(err, apperr.ErrorCodeUnknown, "CreatePostHandler: Handle: repo.FindCategoryByName")
		}

		schemaValidator := validator.NewSchemaValidator(categ.ValueDefinitions)

		err = schemaValidator.Validate(categoryData.Metadata)
		if err != nil {
			validationErrors := &validator.ValidationError{}
			if errors.As(err, &validationErrors) {
				ok = false
				verrs[categoryData.Name] = *validationErrors
			} else {
				return err
			}
		}
	}

	if !ok {
		return &post.PostMetadataError{Errors: verrs}
	}

	return nil
}
