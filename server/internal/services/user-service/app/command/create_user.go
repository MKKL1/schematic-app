package command

import (
	"context"
	"errors"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	"github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domainErr"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/dto"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type CreateUserParams struct {
	Username string
	Sub      uuid.UUID
}

type CreateUserHandler decorator.CommandHandler[CreateUserParams]

type createUserHandler struct {
	repo   postgres.UserRepository
	idNode *snowflake.Node
}

func NewCreateUserHandler(repo postgres.UserRepository, idNode *snowflake.Node) CreateUserHandler {
	return createUserHandler{repo, idNode}
}

func (h createUserHandler) Handle(ctx context.Context, params CreateUserParams) error {
	user := dto.User{
		ID:      dto.UserID(h.idNode.Generate().Int64()),
		Name:    params.Username,
		OidcSub: params.Sub,
	}

	_, err := h.repo.CreateUser(ctx, user)
	if err != nil {
		var e *db.UniqueConstraintError
		if errors.As(err, &e) {
			switch e.Field {
			case "OidcSub":
				return appErr.WrapErrorf(err, domainErr.ErrorCodeSubConflict, "repo.CreateUser")
			case "Name":
				return appErr.WrapErrorf(err, domainErr.ErrorCodeNameConflict, "repo.CreateUser")
			}
		}
		return appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "repo.CreateUser")
	}

	return nil
}
