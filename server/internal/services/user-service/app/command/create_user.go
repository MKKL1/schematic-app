package command

import (
	"context"
	"errors"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	"github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type CreateUserParams struct {
	Username string
	Sub      uuid.UUID
}

type CreateUserHandler decorator.CommandHandler[CreateUserParams]

type createUserHandler struct {
	repo   user.Repository
	idNode *snowflake.Node
}

func NewCreateUserHandler(repo user.Repository, idNode *snowflake.Node) CreateUserHandler {
	return createUserHandler{repo, idNode}
}

func (h createUserHandler) Handle(ctx context.Context, params CreateUserParams) error {
	newUser := user.Model{
		ID:      h.idNode.Generate().Int64(),
		Name:    params.Username,
		OidcSub: params.Sub,
	}

	_, err := h.repo.CreateUser(ctx, newUser)
	if err != nil {
		var e *db.UniqueConstraintError
		if errors.As(err, &e) {
			switch e.Field {
			case "OidcSub":
				return appErr.WrapErrorf(err, user.ErrorCodeSubConflict, "repo.CreateUser")
			case "Name":
				return appErr.WrapErrorf(err, user.ErrorCodeNameConflict, "repo.CreateUser")
			}
		}
		return appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "repo.CreateUser")
	}

	return nil
}
