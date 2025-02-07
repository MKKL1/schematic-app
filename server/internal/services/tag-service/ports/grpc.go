package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/app/command"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	genproto.UnimplementedTagServiceServer
	app app.Application
}

func NewGrpcServer(app app.Application) *GrpcServer {
	return &GrpcServer{app: app}
}

func (g GrpcServer) CreateCategoryVars(ctx context.Context, params *genproto.CreateCategoryVarsParams) (*emptypb.Empty, error) {
	cmd := command.CreateCategoryVarsParams{
		PostId:   params.PostId,
		Category: params.Category,
		Values:   params.Values,
	}

	_, err := g.app.Commands.CreateCategoryVars.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g GrpcServer) GetCategVarsByPost(ctx context.Context, params *genproto.GetCategVarsByPostRequest) (*genproto.CategVarsByPostResponse, error) {
	vars, err := g.app.Queries.GetCategVarsByPost.Handle(ctx, params.GetPostId())
	if err != nil {
		return nil, err
	}

	protoVars := make([]*genproto.CategoryVars, len(vars))
	for i, categoryVars := range vars {
		protoVars[i] = &genproto.CategoryVars{
			PostId:   categoryVars.PostID,
			Category: categoryVars.Category,
			Vars:     categoryVars.Values,
		}
	}

	return &genproto.CategVarsByPostResponse{
		Items: protoVars,
	}, nil
}
