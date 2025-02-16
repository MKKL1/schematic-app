package user

import (
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/http"
	"github.com/labstack/echo/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func RegisterErrorMappers() {
	http.Mapper.AddMapper("USER_NOT_FOUND", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*http.ErrorResponse, bool) {
		return &http.ErrorResponse{
			Errors: []http.ErrorDetail{
				{
					Reason:  info.Reason,
					Message: fmt.Sprintf("User by id %s not found", info.Metadata["id"]),
				},
			},
		}, true
	})

	http.Mapper.AddMapper("USER_NAME_CONFLICT", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*http.ErrorResponse, bool) {
		return &http.ErrorResponse{
			Errors: []http.ErrorDetail{
				{
					Reason:  info.Reason,
					Message: "User by given username already exists",
				},
			},
		}, true
	})

	http.Mapper.AddMapper("USER_SUB_CONFLICT", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*http.ErrorResponse, bool) {
		return &http.ErrorResponse{
			Errors: []http.ErrorDetail{
				{
					Reason:  info.Reason,
					Message: "User by given OIDC subject already registered",
				},
			},
		}, true
	})
}
