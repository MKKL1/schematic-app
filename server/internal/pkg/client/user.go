package client

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/google/uuid"
)

type User struct {
	ID      int64
	Name    string
	OidcSub uuid.UUID
}

func protoToDto(prUser *genproto.User) (*User, error) {
	if prUser == nil {
		return nil, nil
	}

	sub, err := uuid.FromBytes(prUser.OidcSub)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:      prUser.Id,
		Name:    prUser.Name,
		OidcSub: sub,
	}, nil
}

type UserApplication struct {
	Command UserCommandService
	Query   UserQueryService
}

type UserCommandService interface {
	CreateUser(ctx context.Context, name string, sub uuid.UUID) (int64, error)
}

type UserQueryService interface {
	GetUserById(ctx context.Context, id int64) (*User, error)
	GetUserByName(ctx context.Context, name string) (*User, error)
	GetUserBySub(ctx context.Context, sub uuid.UUID) (*User, error)
}

type UserQueryGrpcService struct {
	userServiceClient genproto.UserServiceClient
}

func (u UserQueryGrpcService) GetUserById(ctx context.Context, id int64) (*User, error) {
	byId, err := u.userServiceClient.GetUserById(ctx, &genproto.GetUserByIdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return protoToDto(byId)
}

func (u UserQueryGrpcService) GetUserByName(ctx context.Context, name string) (*User, error) {
	byName, err := u.userServiceClient.GetUserByName(ctx, &genproto.GetUserByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return protoToDto(byName)
}

func (u UserQueryGrpcService) GetUserBySub(ctx context.Context, sub uuid.UUID) (*User, error) {
	subBytes, err := sub.MarshalBinary()
	if err != nil {
		return nil, err
	}

	bySub, err := u.userServiceClient.GetUserBySub(ctx, &genproto.GetUserBySubRequest{
		OidcSub: subBytes,
	})
	if err != nil {
		return nil, err
	}

	return protoToDto(bySub)
}

type UserCommandGrpcService struct {
	userServiceClient genproto.UserServiceClient
}

func (u UserCommandGrpcService) CreateUser(ctx context.Context, name string, sub uuid.UUID) (int64, error) {
	subBytes, err := sub.MarshalBinary()
	if err != nil {
		return 0, err
	}

	newId, err := u.userServiceClient.CreateUser(ctx, &genproto.CreateUserRequest{
		Name:    name,
		OidcSub: subBytes,
	})
	if err != nil {
		return 0, err
	}

	return newId.Id, nil
}

func NewUsersClient(ctx context.Context, addr string) UserApplication {
	conn := NewConnection(ctx, addr)

	service := genproto.NewUserServiceClient(conn)
	query := UserQueryGrpcService{userServiceClient: service}
	command := UserCommandGrpcService{userServiceClient: service}

	return UserApplication{
		Query:   query,
		Command: command,
	}
}
