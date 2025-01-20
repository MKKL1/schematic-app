package query

import (
	"context"
	"errors"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	"github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domainErr"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/dto"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres"
	"github.com/google/uuid"
)

type GetUserBySubParams struct {
	Sub uuid.UUID
}

type GetUserBySubHandler decorator.QueryHandler[GetUserBySubParams, dto.User]

type getUserBySubHandler struct {
	repo postgres.UserRepository
}

func NewGetUserBySubHandler(repo postgres.UserRepository) GetUserBySubHandler {
	return getUserBySubHandler{repo}
}

func (h getUserBySubHandler) Handle(ctx context.Context, params GetUserBySubParams) (dto.User, error) {
	user, err := h.repo.FindByOidcSub(ctx, params.Sub)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return dto.User{}, appErr.WrapErrorf(err, domainErr.ErrorCodeUserNotFound, "repo.FindByOidcSub")
		}
		return dto.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "repo.FindByOidcSub")
	}

	model, err := dto.ToDTO(user)
	if err != nil {
		return dto.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "User.ToDTO")
	}
	return model, nil
}
