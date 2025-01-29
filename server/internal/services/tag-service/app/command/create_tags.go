package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	"github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/domain/tag"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type CreateTagParams struct {
	postId   int64
	tagName  string
	tagValue *string
}

type CreateTagHandler decorator.CommandHandler[CreateTagParams, _]

type createTagHandler struct {
	repo tag.Repository
}

func NewCreateTagHandler(repo tag.Repository) CreateTagHandler {
	return createTagHandler{repo}
}

func (h createTagHandler) Handle(ctx context.Context, params CreateTagParams) error {
	tagEntity := tag.Entity{
		PostID:   params.postId,
		TagName:  params.tagName,
		TagValue: params.tagValue,
	}

	h.repo.AddTag(ctx)
}
