package query

import (
	"context"
	"errors"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/google/uuid"
)

type GetUserBySubParams struct {
	Sub uuid.UUID
}

type GetUserBySubHandler decorator.QueryHandler[GetUserBySubParams, user.User]

type getUserBySubHandler struct {
	repo user.Repository
}

func NewGetUserBySubHandler(repo user.Repository) GetUserBySubHandler {
	return getUserBySubHandler{repo}
}

func (h getUserBySubHandler) Handle(ctx context.Context, params GetUserBySubParams) (user.User, error) {
	userModel, err := h.repo.FindByOidcSub(ctx, params.Sub)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return user.User{}, apperr.WrapErrorf(err, user.ErrorCodeUserNotFound, "repo.FindByOidcSub: user not found by sub: %s", params.Sub.String()).
				AddMetadata("sub", params.Sub.String())
		}
		return user.User{}, apperr.WrapErrorf(err, apperr.ErrorCodeUnknown, "repo.FindByOidcSub: by sub: %s", params.Sub.String()).
			AddMetadata("sub", params.Sub.String())
	}

	return user.EntityToDTO(userModel), nil
}
