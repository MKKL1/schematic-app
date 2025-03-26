package grpc

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DefaultErrorMapper struct {
	Mappers           []func(err error) (error, bool)
	CustomSlugMappers map[string]func(err apperr.SlugError) error
	FallBackMapper    func(err error) error
}

func BuildGrpcError(err apperr.SlugError, grpcCode codes.Code, message string) error {
	st := status.New(grpcCode, message)

	errorInfo := &errdetails.ErrorInfo{
		Reason:   err.Slug,
		Domain:   "schem",
		Metadata: err.Metadata,
	}
	stWithDetails, errDetails := st.WithDetails(errorInfo)
	if errDetails != nil {
		// In the rare case that attaching details fails, log and return the original status.
		return st.Err()
	}
	return stWithDetails.Err()
}

func fallbackMapper(err error) error {
	if slugErr, ok := apperr.FromError(err); ok {
		var code codes.Code
		switch slugErr.Code {
		case apperr.ErrorCodeNotFound:
			code = codes.NotFound
		case apperr.ErrorCodeBadRequest:
			code = codes.InvalidArgument
		case apperr.ErrorCodeConflict:
			code = codes.AlreadyExists
		case apperr.ErrorCodeUnauthorized:
			code = codes.Unauthenticated
		default:
			code = codes.Unknown
		}

		return BuildGrpcError(*slugErr, code, slugErr.Error())
	}

	return status.Error(codes.Internal, err.Error())
}

func NewDefaultErrorMapper() *DefaultErrorMapper {
	return &DefaultErrorMapper{
		Mappers:           make([]func(err error) (error, bool), 0),
		CustomSlugMappers: make(map[string]func(err apperr.SlugError) error),
		FallBackMapper:    fallbackMapper,
	}
}

func (dmm DefaultErrorMapper) Map(err error) error {
	if err == nil {
		return nil
	}

	for _, handler := range dmm.Mappers {
		newErr, ok := handler(err)
		if ok {
			return newErr
		}
	}

	if slugErr, ok := apperr.FromError(err); ok {
		if errMapper, found := dmm.CustomSlugMappers[slugErr.Slug]; found {
			return errMapper(*slugErr)
		}
	}

	return dmm.FallBackMapper(err)
}

func ErrorMappingUnaryInterceptor(mapper func(err error) error) grpc.UnaryServerInterceptor {
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
