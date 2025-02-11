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

type CategoryApplication struct {
	Command CategoryCommandService
	Query   CategoryQueryService
}

type CategoryCommandService interface {
	CreateCategoryVars(ctx context.Context, params CreateCategoryVarsParams) error
}

type CategoryQueryService interface {
	GetCategVarsByPost(ctx context.Context, id int64) ([]PostCategoryVars, error)
}

type CategoryCommandGrpcService struct {
	grpcClient genproto.CategoryServiceClient
}

type CategoryQueryGrpcService struct {
	grpcClient genproto.CategoryServiceClient
}

func (t CategoryCommandGrpcService) CreateCategoryVars(ctx context.Context, params CreateCategoryVarsParams) error {
	_, err := t.grpcClient.CreateCategoryVars(context.Background(), &genproto.CreateCategoryVarsParams{
		PostId:   params.PostId,
		Category: params.Category,
		Values:   params.Values,
	})
	return err
}

func (t CategoryQueryGrpcService) GetCategVarsByPost(ctx context.Context, id int64) ([]PostCategoryVars, error) {
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

func NewCategoryClient(ctx context.Context, addr string) CategoryApplication {
	conn := NewConnection(ctx, addr)

	service := genproto.NewCategoryServiceClient(conn)
	query := CategoryQueryGrpcService{
		grpcClient: service,
	}
	command := CategoryCommandGrpcService{
		grpcClient: service,
	}

	return CategoryApplication{
		Query:   query,
		Command: command,
	}
}
