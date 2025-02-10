package query

import (
	"context"
	"errors"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
)

type GetUserByIdParams struct {
	Id user.ID
}

type GetUserByIdHandler decorator.QueryHandler[GetUserByIdParams, user.User]

type getUserByIdHandler struct {
	repo user.Repository
}

func NewGetUserByIdHandler(repo user.Repository) GetUserByIdHandler {
	return getUserByIdHandler{repo}
}

func (h getUserByIdHandler) Handle(ctx context.Context, params GetUserByIdParams) (user.User, error) {
	userModel, err := h.repo.FindById(ctx, params.Id.Unwrap())
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {

			return user.User{}, apperr.WrapErrorf(err, user.ErrorCodeUserNotFound, "repo.FindById: user not found by id: %s", params.Id.String())
		}
		return user.User{}, apperr.WrapErrorf(err, apperr.ErrorCodeUnknown, "repo.FindById: by id: %s", params.Id.String())
	}

	return user.EntityToDTO(userModel), nil
}
