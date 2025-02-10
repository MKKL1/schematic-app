package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/app/command"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func ErrorMapper(err error) error {
	var code codes.Code
	if appErr, ok := apperr.FromError(err); ok {
		switch appErr.Code {
		default:
			code = codes.Unknown
		}

		st := status.New(code, appErr.Message)

		errorInfo := &errdetails.ErrorInfo{
			Reason:   string(appErr.Code),
			Domain:   "schem.tag",
			Metadata: appErr.Metadata,
		}
		stWithDetails, errDetails := st.WithDetails(errorInfo)
		if errDetails != nil {
			// In the rare case that attaching details fails, log and return the original status.
			return st.Err()
		}
		return stWithDetails.Err()
	}

	st := status.New(codes.Internal, "Internal Server Error")
	return st.Err()
}
