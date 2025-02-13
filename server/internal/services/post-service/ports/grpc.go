package ports

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	vars := make([]*genproto.CategoryVars, len(dto.CategoryVars))
	for i, v := range dto.CategoryVars {
		marshal, err := json.Marshal(v.Values)
		if err != nil {
			log.Fatal().Err(err).Msg("Error marshalling values")
			return nil
		}
		vars[i] = &genproto.CategoryVars{
			Name:     v.CategoryName,
			Metadata: marshal,
		}
	}

	tags := make([]*genproto.Tag, len(dto.Tags))
	for i, v := range dto.Tags {
		tags[i] = &genproto.Tag{
			Tag: v,
		}
	}

	return &genproto.Post{
		Id:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Owner:       dto.Owner,
		Author:      dto.AuthorID,
		Vars:        vars,
		Tags:        tags,
	}
}
