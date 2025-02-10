package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorMapper is a function that converts an error from the domain layer into a gRPC status error.
type ErrorMapper func(err error) error

// DefaultErrorMapper is a fallback that maps any error to codes.Unknown.
func DefaultErrorMapper(err error) error {
	return status.Errorf(codes.Unknown, err.Error())
}

func ErrorMappingUnaryInterceptor(mapper ErrorMapper) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Call the actual handler.
		resp, err := handler(ctx, req)
		if err != nil {
			// If the error isn’t already a gRPC status error, convert it.
			if _, ok := status.FromError(err); !ok {
				err = mapper(err)
			}
		}
		return resp, err
	}
}
