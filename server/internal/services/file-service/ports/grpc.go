package ports

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	genproto.UnimplementedFileServiceServer
	app app.Application
}

func NewGrpcServer(app app.Application) *GrpcServer {
	return &GrpcServer{app: app}
}

func (g GrpcServer) UploadTempFile(stream grpc.ClientStreamingServer[genproto.UploadTempRequest, emptypb.Empty]) error {
	//TODO to stream
}
