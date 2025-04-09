package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/rs/zerolog"
	"strconv"
)

type GetPostByIdParams struct {
	Id int64
}

type GetPostByIdHandler decorator.QueryHandler[GetPostByIdParams, post.Post]

type getPostByIdHandler struct {
	repo post.Repository
}

func NewGetPostByIdHandler(repo post.Repository, logger zerolog.Logger, metrics metrics.Client) GetPostByIdHandler {
	return decorator.ApplyQueryDecorators(
		getPostByIdHandler{repo},
		logger,
		metrics,
	)
}

func (h getPostByIdHandler) Handle(ctx context.Context, params GetPostByIdParams) (post.Post, error) {
	postModel, err := h.repo.FindById(ctx, params.Id)
	if err != nil {
		wrappedErr := fmt.Errorf("finding post by id %d: %w", params.Id, err)
		if errors.Is(err, db.ErrNoRows) {
			return post.Post{}, apperr.NewSlugErrorWithCode(wrappedErr, post.ErrorSlugPostNotFound, apperr.ErrorCodeNotFound).
				AddMetadata("id", strconv.FormatInt(params.Id, 10))
		}
		return post.Post{}, wrappedErr
	}

	return postModel, nil
}
