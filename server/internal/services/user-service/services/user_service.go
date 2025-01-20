package services

import (
	"context"
	"errors"
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	"github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/dto"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository    postgres.UserRepository
	userSnowflakeNode *snowflake.Node
}

func NewUserService(userRepository postgres.UserRepository) *UserService {
	if userRepository == nil {
		panic("userRepository cannot be nil")
	}
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	return &UserService{
		userRepository:    userRepository,
		userSnowflakeNode: node,
	}
}

func (us *UserService) GetUserByOidcSub(ctx context.Context, sub uuid.UUID) (dto.User, error) {
	user, err := us.userRepository.FindByOidcSub(ctx, sub)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return dto.User{}, appErr.WrapErrorf(err, ErrorCodeUserNotFound, "repo.FindByOidcSub")
		}
		return dto.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "repo.FindByOidcSub")
	}

	model, err := dto.ToDTO(user)
	if err != nil {
		return dto.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "User.ToDTO")
	}
	return model, nil
}

func (us *UserService) GetUserById(ctx context.Context, id dto.UserID) (dto.User, error) {
	user, err := us.userRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return dto.User{}, appErr.WrapErrorf(err, ErrorCodeUserNotFound, "repo.FindById")
		}
		return dto.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "repo.FindById")
	}

	model, err := dto.ToDTO(user)
	if err != nil {
		return dto.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "User.ToDTO")
	}
	return model, nil
}

func (us *UserService) CreateUser(ctx context.Context, username string, sub uuid.UUID) (dto.User, error) {

	user := dto.User{
		ID:      dto.UserID(us.userSnowflakeNode.Generate().Int64()),
		Name:    username,
		OidcSub: sub,
	}

	_, err := us.userRepository.CreateUser(ctx, user)
	if err != nil {
		var e *db.UniqueConstraintError
		if errors.As(err, &e) {
			switch e.Field {
			case "OidcSub":
				return dto.User{}, appErr.WrapErrorf(err, ErrorCodeSubConflict, "repo.CreateUser")
			case "Name":
				return dto.User{}, appErr.WrapErrorf(err, ErrorCodeNameConflict, "repo.CreateUser")
			}
		}
		return dto.User{}, appErr.WrapErrorf(err, appErr.ErrorCodeUnknown, "repo.CreateUser")
	}

	return user, nil
}
