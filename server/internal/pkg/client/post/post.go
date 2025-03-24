package post

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
)

type PostApplication struct {
	Command PostCommandService
	Query   PostQueryService
}

type PostCommandService interface {
	CreatePost(ctx context.Context, params CreatePostRequest) (int64, error)
}

type PostQueryService interface {
	GetPostById(ctx context.Context, id int64) (PostDto, error)
}

type PostCommandGrpcService struct {
	grpcClient genproto.PostServiceClient
}

func NewPostCommandGrpcService(grpcClient genproto.PostServiceClient) *PostCommandGrpcService {
	return &PostCommandGrpcService{grpcClient: grpcClient}
}

func (p PostCommandGrpcService) CreatePost(ctx context.Context, params CreatePostRequest) (int64, error) {
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

type PostQueryGrpcService struct {
	grpcClient genproto.PostServiceClient
}

func NewPostQueryGrpcService(grpcClient genproto.PostServiceClient) *PostQueryGrpcService {
	return &PostQueryGrpcService{grpcClient: grpcClient}
}

func (p PostQueryGrpcService) GetPostById(ctx context.Context, id int64) (PostDto, error) {
	post, err := p.grpcClient.GetPostById(ctx, &genproto.PostByIdRequest{Id: id})
	if err != nil {
		return PostDto{}, err
	}

	return ProtoToDto(post)
}

func NewPostClient(ctx context.Context, addr string) PostApplication {
	conn := client.NewConnection(ctx, addr)

	service := genproto.NewPostServiceClient(conn)
	query := PostQueryGrpcService{
		grpcClient: service,
	}
	command := PostCommandGrpcService{
		grpcClient: service,
	}

	return PostApplication{
		Query:   query,
		Command: command,
	}
}
