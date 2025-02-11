package client

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
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
	conn := NewConnection(ctx, addr)

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
