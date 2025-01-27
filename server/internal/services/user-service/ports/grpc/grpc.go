package grpc

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/query"
	domainUser "github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/google/uuid"
)

type GrpcServer struct {
	genproto.UnimplementedUserServiceServer
	app app.Application
}

func NewGrpcServer(app app.Application) *GrpcServer {
	return &GrpcServer{app: app}
}

func (s GrpcServer) GetUserById(ctx context.Context, request *genproto.GetUserByIdRequest) (*genproto.User, error) {
	userDto, err := s.app.Queries.GetUserById.Handle(ctx, query.GetUserByIdParams{Id: domainUser.ID(request.GetId())})
	if err != nil {
		return nil, err
	}

	return dtoToProto(userDto)
}

func (s GrpcServer) GetUserBySub(ctx context.Context, request *genproto.GetUserBySubRequest) (*genproto.User, error) {
	sub, err := uuid.FromBytes(request.GetOidcSub())
	if err != nil {
		return nil, err
	}

	userDto, err := s.app.Queries.GetUserBySub.Handle(ctx, query.GetUserBySubParams{Sub: sub})
	if err != nil {
		return nil, err
	}

	return dtoToProto(userDto)
}

func (s GrpcServer) GetUserByName(ctx context.Context, request *genproto.GetUserByNameRequest) (*genproto.User, error) {
	panic("implement me")
}

func (s GrpcServer) mustEmbedUnimplementedUserServiceServer() {
	panic("implement me")
}

func dtoToProto(userDto domainUser.User) (*genproto.User, error) {
	subBytes, err := userDto.OidcSub.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &genproto.User{
		Id:      userDto.ID.Unwrap(),
		Name:    userDto.Name,
		OidcSub: subBytes,
	}, nil
}
