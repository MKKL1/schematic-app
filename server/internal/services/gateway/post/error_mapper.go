package post

import (
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/grpc"
	gatewayHttp "github.com/MKKL1/schematic-app/server/internal/services/gateway/http"
	"github.com/labstack/echo/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	netHttp "net/http"
)

func RegisterErrorMappers() {
	gatewayHttp.Mapper.AddMapper("POST_NOT_FOUND", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*gatewayHttp.GatewayError, bool) {
		return &gatewayHttp.GatewayError{
			HttpCode: netHttp.StatusNotFound,
			ErrResponse: gatewayHttp.ErrorResponse{
				Errors: []gatewayHttp.ErrorDetail{
					{
						Reason:  info.Reason,
						Message: fmt.Sprintf("Post by id '%s' not found", info.Metadata["id"]),
						Metadata: map[string]string{
							"id": info.Metadata["id"],
						},
					},
				},
			},
		}, true
	})

	gatewayHttp.Mapper.AddMapper("CATEGORY_NOT_FOUND", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*gatewayHttp.GatewayError, bool) {
		return &gatewayHttp.GatewayError{
			HttpCode: netHttp.StatusNotFound,
			ErrResponse: gatewayHttp.ErrorResponse{
				Errors: []gatewayHttp.ErrorDetail{
					{
						Reason:  info.Reason,
						Message: fmt.Sprintf("Category by name '%s' not found", info.Metadata["name"]),
						Metadata: map[string]string{
							"name": info.Metadata["name"],
						},
					},
				},
			},
		}, true
	})
}

func HandlePostCreateErr(err error, requestData PostCreateRequest) error {
	st, ok := status.FromError(err)
	if ok {
		errInfo, found := grpc.GetMessage[errdetails.ErrorInfo](st.Details())
		if !found {
			return err
		}

		if errInfo.GetReason() != "POST_METADATA_VALIDATION_ERROR" {
			return err
		}

		badRequest, found := grpc.GetMessage[errdetails.BadRequest](st.Details())
		if !found {
			return err
		}

		var errDetails []gatewayHttp.ErrorDetail
		for _, v := range badRequest.GetFieldViolations() {
			parameter := mapFieldPath(v.GetField(), requestData)
			errDetails = append(errDetails, gatewayHttp.ValidationErrorBuilder{
				Parameter: parameter,
				Detail:    v.GetReason(),
				Message:   v.GetDescription(),
			}.Build())
		}

		return &gatewayHttp.GatewayError{
			HttpCode: netHttp.StatusBadRequest,
			ErrResponse: gatewayHttp.ErrorResponse{
				Errors: errDetails,
			},
		}
	}

	return err
}
