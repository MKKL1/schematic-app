package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/http"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"strconv"
)

type CreatePostParams struct {
	http.PostCreateRequest
	Sub uuid.UUID
}

type CreatePostHandler decorator.CommandHandler[CreatePostParams, int64]

type createPostHandler struct {
	repo        post.Repository
	idNode      *snowflake.Node
	userService client.UserGrpcService
}

func NewCreatePostHandler(repo post.Repository, idNode *snowflake.Node, userService client.UserGrpcService) CreatePostHandler {
	return createPostHandler{repo, idNode, userService}
}

func (h createPostHandler) Handle(ctx context.Context, params CreatePostParams) (int64, error) {
	user, err := h.userService.GetUserBySub(ctx, params.Sub)
	if err != nil {
		return 0, err
	}

	var authorName *string
	var authorID *int64
	if params.Author != nil {
		authorName = params.Author.Name

		if params.Author.ID != nil {
			parInt, err := strconv.ParseInt(*params.Author.ID, 10, 64)
			if err != nil {
				return 0, err
			}
			authorID = &parInt
		}
	}

	newPost := post.Model{
		ID:          h.idNode.Generate().Int64(),
		Name:        params.Name,
		Description: params.Description,
		Owner:       user.ID,
		AuthorName:  authorName,
		AuthorID:    authorID,
	}

	err = h.repo.Create(ctx, newPost)
	if err != nil {
		return 0, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "repo.Create")
	}

	return newPost.ID, nil
}
