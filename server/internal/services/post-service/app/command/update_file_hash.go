package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UpdateFileHashCommand struct {
	Ids []FileHashTempId
}

type FileHashTempId struct {
	TempId uuid.UUID
	Hash   string
}

type UpdateFileHashHandler decorator.CommandHandler[UpdateFileHashCommand, any]

type updateFileHashHandler struct {
	repo post.Repository
}

func NewUpdateAttachedFilesHandler(repo post.Repository, logger zerolog.Logger, metrics metrics.Client) UpdateFileHashHandler {
	return decorator.ApplyCommandDecorators(
		updateFileHashHandler{repo: repo},
		logger,
		metrics,
	)
}

func (u updateFileHashHandler) Handle(ctx context.Context, cmd UpdateFileHashCommand) (any, error) {
	err := u.repo.UpdateFileHashByTempId(ctx, updateFileHashCmdToDomain(cmd))
	return err, nil
}

func updateFileHashCmdToDomain(cmd UpdateFileHashCommand) []post.UpdateFileHashParams {
	params := make([]post.UpdateFileHashParams, len(cmd.Ids))
	for i, id := range cmd.Ids {
		params[i] = post.UpdateFileHashParams{
			Hash:   id.Hash,
			TempId: id.TempId,
		}
	}
	return params
}
