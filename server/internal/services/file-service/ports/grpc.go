package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewFileGrpcErrorMapper() func(error) error {
	mapper := grpc.NewDefaultErrorMapper()
	return mapper.Map
}

type GrpcServer struct {
	genproto.UnimplementedFileServiceServer
	app app.Application
}

func NewGrpcServer(app app.Application) *GrpcServer {
	return &GrpcServer{app: app}
}

func (g GrpcServer) DeleteExpiredFiles(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	_, err := g.app.Commands.DeleteExpiredFiles.Handle(ctx, command.DeleteExpiredFilesParams{})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
