package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client/user"
	"github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category/validator"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type CreatePost struct {
	Name        string
	Description *string
	AuthorID    *int64
	Sub         uuid.UUID
	Categories  []CreatePostCategory
	Tags        []string
	Files       []uuid.UUID
}

type CreatePostCategory struct {
	Name     string
	Metadata map[string]interface{}
}

type CreatePostHandler decorator.CommandHandler[CreatePost, int64]

type createPostHandler struct {
	repo         post.Repository
	categoryRepo category.Repository
	idNode       *snowflake.Node
	userService  user.Service
	eventBus     *cqrs.EventBus
	logger       zerolog.Logger
}

func NewCreatePostHandler(repo post.Repository, categoryRepo category.Repository, idNode *snowflake.Node, userService user.Service, eventBus *cqrs.EventBus, logger zerolog.Logger, metrics metrics.Client) CreatePostHandler {
	return decorator.ApplyCommandDecorators(
		createPostHandler{repo, categoryRepo, idNode, userService, eventBus, logger},
		logger,
		metrics,
	)
}

func (h createPostHandler) Handle(ctx context.Context, params CreatePost) (int64, error) {
	logger := decorator.AddCmdInfo(params, h.logger)
	err := h.validateCategories(ctx, params.Categories) //TODO validation should be in domain
	if err != nil {
		return 0, fmt.Errorf("validating post categories: %w", err)
	}

	usr, err := h.userService.GetUserBySub(ctx, params.Sub)
	if err != nil {
		return 0, fmt.Errorf("querying user by sub: %w", err)
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
			Name:   "TODO", //TODO get name from file service, this should be done when checking if file exists
			TempId: f.String(),
		}
	}

	newPost := post.CreatePostParams{ //TODO construct in domain
		ID:          h.idNode.Generate().Int64(),
		Name:        params.Name,
		Description: params.Description,
		AuthorID:    params.AuthorID,
		Owner:       usr.ID,
		Categories:  categs,
		Tags:        params.Tags,
		Files:       files,
	}

	err = h.repo.Create(ctx, newPost)
	if err != nil {
		return 0, fmt.Errorf("creating post in repo: %w", err)
	}

	go func() {
		err = h.publishCreatePostEvent(ctx, newPost)
		if err != nil {
			//TODO handle, log or delete created post from database
			logger.Error().Err(err).Msg("Failed to publish create post event")
		}
	}()

	return newPost.ID, nil
}

func (h createPostHandler) publishCreatePostEvent(ctx context.Context, newPost post.CreatePostParams) error {
	categs := make(post.PostCategories, len(newPost.Categories))
	for _, c := range newPost.Categories {
		categs[c.Name] = c.Metadata
	}

	eventFiles := make([]post.PostCreatedFileData, len(newPost.Files))
	for i, f := range newPost.Files {
		eventFiles[i] = post.PostCreatedFileData{
			TempId: f.TempId,
			Name:   f.Name,
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
		return fmt.Errorf("publishing event: %w", err)
	}
	return nil
}

// TODO may be in domain??
func (h createPostHandler) validateCategories(ctx context.Context, categories []CreatePostCategory) error {
	verrs := make(map[string]validator.ValidationError)
	var ok = true
	for _, categoryData := range categories {
		//TODO can be bulk
		categoryModel, err := h.categoryRepo.FindCategoryByName(ctx, categoryData.Name)
		if err != nil {
			wrappedErr := fmt.Errorf("finding category %q: %w", categoryData.Name, err)
			if errors.Is(err, db.ErrNoRows) {
				return apperr.NewSlugErrorWithCode(wrappedErr, category.ErrorSlugCategoryNotFound, apperr.ErrorCodeNotFound).
					AddMetadata("name", categoryData.Name)
			}
			return wrappedErr
		}

		schemaValidator := validator.NewSchemaValidator(categoryModel.MetadataSchema)
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
