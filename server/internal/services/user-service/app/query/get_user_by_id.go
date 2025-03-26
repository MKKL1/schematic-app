package query

import (
	"context"
	"errors"
	"fmt"
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
		wrappedErr := fmt.Errorf("finding user by id %s in repo: %w", params.Id.String(), err)
		if errors.Is(err, db.ErrNoRows) {
			return user.User{}, apperr.NewSlugErrorWithCode(wrappedErr, user.ErrorSlugUserNotFound, apperr.ErrorCodeNotFound)
		}
		return user.User{}, wrappedErr
	}

	return user.EntityToDTO(userModel), nil
}
