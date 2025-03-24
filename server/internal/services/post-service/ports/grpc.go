package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/mappers"
)

type GrpcServer struct {
	genproto.UnimplementedPostServiceServer
	app app.Application
}

func NewGrpcServer(app app.Application) *GrpcServer {
	return &GrpcServer{app: app}
}

func (g GrpcServer) GetPostById(ctx context.Context, request *genproto.PostByIdRequest) (*genproto.Post, error) {
	dto, err := g.app.Queries.GetPostById.Handle(ctx, query.GetPostByIdParams{Id: request.GetId()})
	if err != nil {
		return nil, err
	}

	return mappers.AppToProto(dto)
}

func (g GrpcServer) CreatePost(ctx context.Context, request *genproto.CreatePostRequest) (*genproto.CreatePostResponse, error) {
	cmd, err := mappers.CreatePostRequestProtoToCmd(request)
	if err != nil {
		return nil, err
	}

	createdId, err := g.app.Commands.CreatePost.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &genproto.CreatePostResponse{Id: createdId}, nil
}
