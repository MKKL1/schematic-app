package query

import (
	"context"
	"errors"
	"fmt"
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
		wrappedErr := fmt.Errorf("finding user by sub %s in repo: %w", params.Sub.String(), err)
		if errors.Is(err, db.ErrNoRows) {
			return user.User{}, apperr.NewSlugErrorWithCode(wrappedErr, user.ErrorSlugUserNotFound, apperr.ErrorCodeNotFound).
				AddMetadata("sub", params.Sub.String())
		}
		return user.User{}, wrappedErr
	}

	return user.EntityToDTO(userModel), nil
}
