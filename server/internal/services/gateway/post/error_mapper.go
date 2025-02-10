package post

import (
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/http"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func RegisterErrorMappers() {
	http.Mapper.AddMapper("POST_NOT_FOUND", func(info *errdetails.ErrorInfo) http.ErrorDetail {
		return http.ErrorDetail{
			Domain:  info.Domain,
			Reason:  info.Reason,
			Message: fmt.Sprintf("Post by id %s not found", info.Metadata["id"]),
		}
	})
}
