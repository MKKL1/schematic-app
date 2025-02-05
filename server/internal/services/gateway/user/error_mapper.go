package user

import (
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	httpErr "github.com/MKKL1/schematic-app/server/internal/pkg/http"
	domain "github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"net/http"
)

var (
	ErrUserNotFound     = httpErr.NewErrorResponse("USER_NOT_FOUND", http.StatusNotFound, "User not found")
	ErrUserNameConflict = httpErr.NewErrorResponse("USER_NAME_CONFLICT", http.StatusConflict, "Username already exists")
	ErrUserSubConflict  = httpErr.NewErrorResponse("USER_SUB_CONFLICT", http.StatusConflict, "OIDC subject already registered")
)

func MapAppError(err appErr.Error) error {
	switch err.Code() {
	case domain.ErrCodeUserNotFound:
		return ErrUserNotFound
	case domain.ErrCodeNameConflict:
		return ErrUserNameConflict
	case domain.ErrCodeSubConflict:
		return ErrUserSubConflict
	}
	//TODO log error
	return &err
}
