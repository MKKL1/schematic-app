package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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

	return dtoToProto(dto), nil
}

func (g GrpcServer) CreatePost(ctx context.Context, request *genproto.CreatePostRequest) (*genproto.CreatePostResponse, error) {
	sub, err := uuid.FromBytes(request.AuthSub)
	if err != nil {
		return nil, err
	}

	createdId, err := g.app.Commands.CreatePost.Handle(ctx, command.CreatePostParams{
		Name:        request.GetName(),
		Description: request.Description,
		AuthorName:  request.AuthorName,
		AuthorID:    request.AuthorId,
		Sub:         sub,
	})
	if err != nil {
		return nil, err
	}

	return &genproto.CreatePostResponse{Id: createdId}, nil
}

func dtoToProto(dto post.Post) *genproto.Post {
	p := &genproto.Post{
		Id:    dto.ID,
		Name:  dto.Name,
		Owner: dto.Owner,
	}

	if dto.Description != nil {
		p.Description = proto.String(*dto.Description)
	}
	if dto.AuthorID != nil {
		p.Author = proto.Int64(*dto.AuthorID)
	}

	p.Vars = &anypb.Any{
		TypeUrl: "CategoryVars",
		Value:   dto.CategoryVars,
	}

	p.Tags = make([]*genproto.Tag, len(dto.Tags))
	for i, t := range dto.Tags {
		p.Tags[i] = &genproto.Tag{Tag: t}
	}

	return p
}
