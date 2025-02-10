package http

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Mapper = NewDefaultErrorMessageMapper()

// GRPCErrorToDetailedHTTPResponse converts a gRPC error (with details) to an HTTP status code and a JSON payload.
func GRPCErrorToDetailedHTTPResponse(err error) (int, interface{}) {
	st, ok := status.FromError(err)
	if !ok {
		st = status.New(codes.Unknown, err.Error())
	}
	httpStatus := runtime.HTTPStatusFromCode(st.Code())
	userMessage := st.Message()
	log.Error().Msg(userMessage)

	var errDetail ErrorDetail

	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			errDetail = Mapper.MapError(info)
			break
		}
	}

	return httpStatus, ErrorResponse{
		Errors: []ErrorDetail{errDetail},
	}
}

// GRPCErrorMiddleware is an Echo middleware that intercepts errors and converts gRPC errors into HTTP JSON responses.
func GRPCErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Execute the next handler.
		err := next(c)
		if err != nil {
			// Check if the error is a gRPC error.
			if _, ok := status.FromError(err); ok {
				httpStatus, payload := GRPCErrorToDetailedHTTPResponse(err)
				return c.JSON(httpStatus, payload)
			}
			// For non–gRPC errors, you could choose to return them as is or wrap them.
			return err
		}
		return nil
	}
}
