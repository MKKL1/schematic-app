package post

import (
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	httpErr "github.com/MKKL1/schematic-app/server/internal/pkg/http"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"net/http"
)

var (
	PostNotFoundError = httpErr.NewErrorResponse("POST_NOT_FOUND", http.StatusNotFound, "Post not found")
)

//func MapError(err error) error {
//	var e *appErr.Error
//	if errors.As(err, &e) {
//		return MapAppError(*e)
//	}
//	return err
//}

func MapAppError(err appErr.Error) error {
	switch err.Code() {
	case post.ErrorCodePostNotFound:
		return PostNotFoundError
	}
	//TODO log error
	return &err
}
