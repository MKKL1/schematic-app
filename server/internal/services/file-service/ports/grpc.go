package ports

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type GrpcServer struct {
	genproto.UnimplementedFileServiceServer
	app app.Application
}

func NewGrpcServer(app app.Application) *GrpcServer {
	return &GrpcServer{app: app}
}

func (g GrpcServer) UploadTempFile(stream grpc.ClientStreamingServer[genproto.UploadTempRequest, genproto.UploadTempFileResponse]) error {
	// Expect the first message to be metadata.
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to receive metadata: %v", err)
	}
	metadata := req.GetMetadata()
	if metadata == nil {
		return status.Error(codes.InvalidArgument, "metadata expected as first message")
	}

	// Use file size and filename from metadata.
	fileSize := metadata.Size
	objectName := metadata.FileName // Use the filename provided in metadata

	// Create a pipe to stream file data.
	pr, pw := io.Pipe()
	errCh := make(chan error, 1)

	// Start a goroutine to read the file chunks from the stream.
	go func() {
		// Ensure that the pipe writer is closed at the end.
		defer func() {
			_ = pw.Close() // Ignore error here since it might have been closed with an error already.
		}()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				// Normal termination of the stream.
				errCh <- nil
				return
			}
			if err != nil {
				// If an error occurs, close the writer with the error and forward it.
				_ = pw.CloseWithError(err)
				errCh <- err
				return
			}
			chunk := req.GetData()
			if len(chunk) > 0 {
				if _, err := pw.Write(chunk); err != nil {
					// Write error: close the writer with error and send it on the channel.
					_ = pw.CloseWithError(err)
					errCh <- err
					return
				}
			}
		}
	}()

	// Use the pipe reader to stream the data to MinIO.
	resp, err := g.app.Commands.UploadTempFile.Handle(
		stream.Context(), command.UploadTempFileParams{
			Reader:      pr,
			FileName:    objectName,
			FileSize:    fileSize,
			ContentType: metadata.ContentType,
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "failed to upload file to MinIO: %v", err)
	}
	//fmt.Printf("Successfully uploaded object %s with %d bytes\n", objectName, info)

	// Check if there was any error in the stream goroutine.
	if streamErr := <-errCh; streamErr != nil {
		return status.Errorf(codes.Internal, "error while receiving stream: %v", streamErr)
	}

	return stream.SendAndClose(&genproto.UploadTempFileResponse{
		Key: resp.Key,
		Exp: resp.Expiration.Milliseconds(),
		Url: resp.Url,
	})
}
