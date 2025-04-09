package post

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
)

type Service interface {
	CreatePost(ctx context.Context, params CreatePostRequest) (int64, error)
	GetPostById(ctx context.Context, id int64) (PostDto, error)
}

type GrpcService struct {
	grpcClient genproto.PostServiceClient
}

func NewGrpcService(grpcClient genproto.PostServiceClient) *GrpcService {
	return &GrpcService{grpcClient: grpcClient}
}

func (p GrpcService) CreatePost(ctx context.Context, params CreatePostRequest) (int64, error) {
	protoRequest, err := CreatePostRequestDtoToProto(params)
	if err != nil {
		return 0, err
	}
	createdId, err := p.grpcClient.CreatePost(ctx, protoRequest)
	if err != nil {
		return 0, err
	}

	return createdId.Id, err
}

func (p GrpcService) GetPostById(ctx context.Context, id int64) (PostDto, error) {
	post, err := p.grpcClient.GetPostById(ctx, &genproto.PostByIdRequest{Id: id})
	if err != nil {
		return PostDto{}, err
	}

	return ProtoToDto(post)
}

func NewPostClient(ctx context.Context, addr string) Service {
	conn := grpc.NewClient(ctx, addr)

	service := genproto.NewPostServiceClient(conn)
	return GrpcService{grpcClient: service}
}
