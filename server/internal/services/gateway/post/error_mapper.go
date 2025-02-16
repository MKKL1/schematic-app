package post

import (
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/http"
	"github.com/labstack/echo/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func RegisterErrorMappers() {
	http.Mapper.AddMapper("POST_NOT_FOUND", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*http.ErrorResponse, bool) {
		return &http.ErrorResponse{
			Errors: []http.ErrorDetail{
				{
					Reason:  info.Reason,
					Message: fmt.Sprintf("Post by id %s not found", info.Metadata["id"]),
				},
			},
		}, true
	})

	//http.Mapper.AddMapper("POST_METADATA_VALIDATION_ERROR", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*http.ErrorResponse, bool) {
	//	badRequest, found := grpc.GetMessage[errdetails.BadRequest](details)
	//	if !found {
	//		return nil, false
	//	}
	//
	//	var errDetails []http.ErrorDetail
	//	for _, v := range badRequest.GetFieldViolations() {
	//		errDetails = append(errDetails, http.ValidationErrorBuilder{
	//			Parameter: v.GetField(),
	//			Detail:    v.GetReason(),
	//			Code:      "?",
	//			Value:     "?",
	//			Message:   v.GetDescription(),
	//		}.Build())
	//	}
	//
	//	return &http.ErrorResponse{
	//		Errors: errDetails,
	//	}, true
	//})
}
