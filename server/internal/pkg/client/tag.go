package client

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostCategoryVars struct {
	PostID   int64
	Category string
	Values   json.RawMessage
}

type CreateCategoryVarsParams struct {
	PostId   int64
	Category string
	Values   []byte
}

type TagApplication struct {
	Command TagCommandService
	Query   TagQueryService
}

type TagCommandService interface {
	CreateCategoryVars(ctx context.Context, params CreateCategoryVarsParams) error
}

type TagQueryService interface {
	GetCategVarsByPost(ctx context.Context, id int64) ([]PostCategoryVars, error)
}

type TagCommandGrpcService struct {
	grpcClient genproto.TagServiceClient
}

type TagQueryGrpcService struct {
	grpcClient genproto.TagServiceClient
}

func (t TagCommandGrpcService) CreateCategoryVars(ctx context.Context, params CreateCategoryVarsParams) error {
	_, err := t.grpcClient.CreateCategoryVars(context.Background(), &genproto.CreateCategoryVarsParams{
		PostId:   params.PostId,
		Category: params.Category,
		Values:   params.Values,
	})
	return err
}

func (t TagQueryGrpcService) GetCategVarsByPost(ctx context.Context, id int64) ([]PostCategoryVars, error) {
	vars, err := t.grpcClient.GetCategVarsByPost(ctx, &genproto.GetCategVarsByPostRequest{
		PostId: id,
	})
	if err != nil {
		return nil, err
	}

	dtoVars := make([]PostCategoryVars, len(vars.Items))
	for i, v := range vars.Items {
		dtoVars[i] = PostCategoryVars{
			PostID:   v.PostId,
			Category: v.Category,
			Values:   v.Vars,
		}
	}

	return dtoVars, nil
}

func NewTagClient(ctx context.Context, addr string) TagApplication {
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

	service := genproto.NewTagServiceClient(conn)
	query := TagQueryGrpcService{
		grpcClient: service,
	}
	command := TagCommandGrpcService{
		grpcClient: service,
	}

	return TagApplication{
		Query:   query,
		Command: command,
	}
}
