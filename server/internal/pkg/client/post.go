package client

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Post struct {
	ID          int64
	Name        string
	Description *string
	Owner       int64
	AuthorID    *int64
}

type CreatePostParams struct {
	Name        string
	Description *string
	AuthorName  *string
	AuthorID    *int64
	Sub         uuid.UUID
}

type PostApplication struct {
	Command PostCommandService
	Query   PostQueryService
}

type PostCommandService interface {
	CreatePost(ctx context.Context, params CreatePostParams) (int64, error)
}

type PostQueryService interface {
	GetPostById(ctx context.Context, id int64) (Post, error)
}

type PostCommandGrpcService struct {
	grpcClient genproto.PostServiceClient
}

func NewPostCommandGrpcService(grpcClient genproto.PostServiceClient) *PostCommandGrpcService {
	return &PostCommandGrpcService{grpcClient: grpcClient}
}

func (p PostCommandGrpcService) CreatePost(ctx context.Context, params CreatePostParams) (int64, error) {
	subBytes, err := params.Sub.MarshalBinary()
	if err != nil {
		return 0, err
	}

	createdId, err := p.grpcClient.CreatePost(ctx, &genproto.CreatePostRequest{
		Name:        params.Name,
		Description: params.Description,
		AuthorName:  params.AuthorName,
		AuthorId:    params.AuthorID,
		AuthSub:     subBytes,
	})
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

func (p PostQueryGrpcService) GetPostById(ctx context.Context, id int64) (Post, error) {
	post, err := p.grpcClient.GetPostById(ctx, &genproto.PostByIdRequest{Id: id})
	if err != nil {
		return Post{}, err
	}

	return postProtoToDto(post), nil
}

func postProtoToDto(post *genproto.Post) Post {
	return Post{
		ID:          post.Id,
		Name:        post.Name,
		Description: post.Description,
		Owner:       post.Owner,
		AuthorID:    post.Author,
	}
}

func NewPostClient(ctx context.Context, addr string) PostApplication {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
