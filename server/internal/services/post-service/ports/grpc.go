package ports

import (
	"context"
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/mappers"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewPostGrpcErrorMapper() func(error) error {
	mapper := grpc.NewDefaultErrorMapper()
	mapper.Mappers = append(mapper.Mappers, func(err error) (error, bool) {
		var pme *post.PostMetadataError
		if !errors.As(err, &pme) {
			return nil, false
		}

		br := &errdetails.BadRequest{}

		for categ, v := range pme.Errors {
			for _, k := range v.Errors {
				br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
					Field:       fmt.Sprintf("%s:%s", categ, k.Field),
					Description: k.Message,
				})
			}
		}

		st := status.New(codes.InvalidArgument, "invalid post metadata")
		errorInfo := &errdetails.ErrorInfo{
			Reason:   post.ErrorSlugPostMetadataValidation,
			Domain:   "schem.post",
			Metadata: nil,
		}
		stWithDetails, errDetails := st.WithDetails(errorInfo, br)
		if errDetails != nil {
			return st.Err(), false
		}
		return stWithDetails.Err(), true
	})

	return mapper.Map
}

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
