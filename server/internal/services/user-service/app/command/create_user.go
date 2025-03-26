package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type CreateUserParams struct {
	Username string
	Sub      uuid.UUID
}

type CreateUserHandler decorator.CommandHandler[CreateUserParams, user.ID]

type createUserHandler struct {
	repo   user.Repository
	idNode *snowflake.Node
}

func NewCreateUserHandler(repo user.Repository, idNode *snowflake.Node) CreateUserHandler {
	return createUserHandler{repo, idNode}
}

func (h createUserHandler) Handle(ctx context.Context, params CreateUserParams) (user.ID, error) {
	newUser := user.Entity{
		ID:      h.idNode.Generate().Int64(),
		Name:    params.Username,
		OidcSub: params.Sub,
	}

	_, err := h.repo.CreateUser(ctx, newUser)
	if err != nil {
		//Implementing domain logic with repository //TODO refactor
		wrappedErr := fmt.Errorf("creating user %s in repo: %w", params.Sub.String(), err)

		var e *db.UniqueConstraintError
		if errors.As(err, &e) {
			switch e.Field {
			case "OidcSub":
				return 0, apperr.NewSlugErrorWithCode(wrappedErr, user.ErrorSlugSubConflict, apperr.ErrorCodeConflict)
			case "Name":
				return 0, apperr.NewSlugErrorWithCode(wrappedErr, user.ErrorSlugNameConflict, apperr.ErrorCodeConflict)
			}
		}
		return 0, wrappedErr
	}

	return user.ID(newUser.ID), nil
}
