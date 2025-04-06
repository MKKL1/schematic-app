package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/mappers"
)

func NewImageGrpcErrorMapper() func(error) error {
	mapper := grpc.NewDefaultErrorMapper()
	return mapper.Map
}

// GrpcServer implements the ImageService.
type GrpcServer struct {
	genproto.UnimplementedImageServiceServer
	app app.Application
}

// NewGrpcServer returns a new ImageService gRPC server.
func NewGrpcServer(app app.Application) *GrpcServer {
	return &GrpcServer{app: app}
}

func (g GrpcServer) GetImageSizes(ctx context.Context, req *genproto.GetImageSizesRequest) (*genproto.GetImageSizesResponse, error) {
	// Create the command parameters from the request.
	params := query.GetImageParams{
		ImageID: req.GetHash(),
	}

	// Handle the query.
	result, err := g.app.Queries.GetImageSizes.Handle(ctx, params)
	if err != nil {
		return nil, err
	}

	// Map the result to the gRPC response.
	return mappers.AppToProtoGetImageSizesResponse(result)
}
