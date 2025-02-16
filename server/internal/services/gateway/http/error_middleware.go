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
func GRPCErrorToDetailedHTTPResponse(c echo.Context, err error) (int, interface{}) {
	st, ok := status.FromError(err)
	if !ok {
		st = status.New(codes.Unknown, err.Error())
	}
	httpStatus := runtime.HTTPStatusFromCode(st.Code())
	userMessage := st.Message()
	log.Error().Msg(userMessage)

	var errResp ErrorResponse

	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			mappedErr, ok := Mapper.MapError(c, info, st.Details())
			if !ok || mappedErr == nil {
				errResp = ErrorResponse{
					Errors: []ErrorDetail{
						{
							Reason:   "INTERNAL_ERROR",
							Message:  "An unexpected error occurred",
							Metadata: map[string]string{},
						},
					},
				}
				break
			}

			errResp = *mappedErr
			break
		}
	}

	return httpStatus, errResp
}
