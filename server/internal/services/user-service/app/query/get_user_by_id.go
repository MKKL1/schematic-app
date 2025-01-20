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
)

type GetUserByIdParams struct {
	Id dto.UserID
}

type GetUserByIdHandler decorator.QueryHandler[GetUserByIdParams, dto.User]

type getUserByIdHandler struct {
	repo postgres.UserRepository
}

func NewGetUserByIdHandler(repo postgres.UserRepository) GetUserByIdHandler {
	return getUserByIdHandler{repo}
}

func (h getUserByIdHandler) Handle(ctx context.Context, params GetUserByIdParams) (dto.User, error) {
	user, err := h.repo.FindById(ctx, params.Id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return dto.User{}, appErr.WrapErrorf(err, domainErr.ErrorCodeUserNotFound, "repo.FindById")
		}
		return dto.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "repo.FindById")
	}

	model, err := dto.ToDTO(user)
	if err != nil {
		return dto.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "User.ToDTO")
	}
	return model, nil
}
