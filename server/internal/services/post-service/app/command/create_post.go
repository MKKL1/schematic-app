package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type CreatePostParams struct {
	Name        string
	Description *string
	AuthorName  *string
	AuthorID    *int64
	Sub         uuid.UUID
}

type CreatePostHandler decorator.CommandHandler[CreatePostParams, int64]

type createPostHandler struct {
	repo        post.Repository
	idNode      *snowflake.Node
	userService client.UserApplication
}

func NewCreatePostHandler(repo post.Repository, idNode *snowflake.Node, userService client.UserApplication) CreatePostHandler {
	return createPostHandler{repo, idNode, userService}
}

func (h createPostHandler) Handle(ctx context.Context, params CreatePostParams) (int64, error) {
	user, err := h.userService.Query.GetUserBySub(ctx, params.Sub)
	if err != nil {
		return 0, err
	}

	//TODO if author id null and author name not null, then create author in author service

	newPost := post.Entity{
		ID:          h.idNode.Generate().Int64(),
		Name:        params.Name,
		Description: params.Description,
		Owner:       user.ID,
		AuthorID:    params.AuthorID,
	}

	err = h.repo.Create(ctx, newPost)
	if err != nil {
		return 0, apperr.WrapErrorf(err, apperr.ErrorCodeUnknown, "repo.Create")
	}

	return newPost.ID, nil
}
