package query

import (
	"context"
	"errors"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
)

type GetPostByIdParams struct {
	Id int64
}

type GetPostByIdHandler decorator.QueryHandler[GetPostByIdParams, post.Post]

type getPostByIdHandler struct {
	repo post.Repository
}

func NewGetPostByIdHandler(repo post.Repository) GetPostByIdHandler {
	return getPostByIdHandler{repo}
}

func (h getPostByIdHandler) Handle(ctx context.Context, params GetPostByIdParams) (post.Post, error) {
	postModel, err := h.repo.FindById(ctx, params.Id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return post.Post{}, apperr.WrapErrorf(err, post.ErrorCodePostNotFound, "repo.FindById")
		}
		return post.Post{}, apperr.WrapErrorf(err, apperr.ErrorCodeUnknown, "repo.FindById")
	}

	return post.ToDTO(postModel), nil
}
