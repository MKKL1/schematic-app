package query

import (
	"context"
	"errors"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	"github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
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
			return user.User{}, appErr.WrapErrorf(err, user.ErrorCodeUserNotFound, "repo.FindByOidcSub")
		}
		return user.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "repo.FindByOidcSub")
	}

	return user.ToDTO(userModel), nil
}
