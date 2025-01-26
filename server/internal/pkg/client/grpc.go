package client

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
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

type UserService interface {
	GetUserById(ctx context.Context, id int64) (*User, error)
	GetUserByName(ctx context.Context, name string) (*User, error)
	GetUserBySub(ctx context.Context, sub uuid.UUID) (*User, error)
}

type UserGrpcService struct {
	userServiceClient genproto.UserServiceClient
}

func (u UserGrpcService) GetUserById(ctx context.Context, id int64) (*User, error) {
	byId, err := u.userServiceClient.GetUserById(ctx, &genproto.GetUserByIdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return protoToDto(byId)
}

func (u UserGrpcService) GetUserByName(ctx context.Context, name string) (*User, error) {
	byName, err := u.userServiceClient.GetUserByName(ctx, &genproto.GetUserByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return protoToDto(byName)
}

func (u UserGrpcService) GetUserBySub(ctx context.Context, sub uuid.UUID) (*User, error) {
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

func NewUsersClient(ctx context.Context, addr string) UserGrpcService {
	conn, err := grpc.NewClient(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Str("addr", addr).Msg("shutting down gRPC server")
				err := conn.Close()
				if err != nil {
					log.Error().Str("addr", addr).Err(err).Msg("failed to close gRPC connection")
					return
				}
				log.Info().Msg("server shut down")
				return
			}
		}
	}()

	return UserGrpcService{userServiceClient: genproto.NewUserServiceClient(conn)}
}
