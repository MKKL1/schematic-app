package user

import (
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/http"
	"github.com/labstack/echo/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	netHttp "net/http"
)

func RegisterErrorMappers() {
	http.Mapper.AddMapper("USER_NOT_FOUND", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*http.GatewayError, bool) {
		metadata := map[string]string{
			"id": info.Metadata["id"],
		}
		message := fmt.Sprintf("User by id %s not found", info.Metadata["id"])
		return http.NewSimpleGatewayError(netHttp.StatusNotFound, info.Reason, message, metadata), true
	})

	http.Mapper.AddMapper("USER_NAME_CONFLICT", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*http.GatewayError, bool) {
		return http.NewSimpleGatewayError(netHttp.StatusConflict, info.Reason, "User by given username already exists", nil), true
	})

	http.Mapper.AddMapper("USER_SUB_CONFLICT", func(c echo.Context, info *errdetails.ErrorInfo, details []any) (*http.GatewayError, bool) {
		return http.NewSimpleGatewayError(netHttp.StatusConflict, info.Reason, "User by given OIDC subject already registered", nil), true
	})
}
