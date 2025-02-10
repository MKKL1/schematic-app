package user

import (
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/http"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func RegisterErrorMappers() {
	http.Mapper.AddMapper("USER_NOT_FOUND", func(info *errdetails.ErrorInfo) http.ErrorDetail {
		return http.ErrorDetail{
			Domain:  info.Domain,
			Reason:  info.Reason,
			Message: fmt.Sprintf("User by id %s not found", info.Metadata["id"]),
		}
	})

	http.Mapper.AddMapper("USER_NAME_CONFLICT", func(info *errdetails.ErrorInfo) http.ErrorDetail {
		return http.ErrorDetail{
			Domain:  info.Domain,
			Reason:  info.Reason,
			Message: "User by given username already exists",
		}
	})

	http.Mapper.AddMapper("USER_SUB_CONFLICT", func(info *errdetails.ErrorInfo) http.ErrorDetail {
		return http.ErrorDetail{
			Domain:  info.Domain,
			Reason:  info.Reason,
			Message: "User by given OIDC subject already registered",
		}
	})
}
