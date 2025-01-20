package http

import (
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	httpErr "github.com/MKKL1/schematic-app/server/internal/pkg/http"
	error2 "github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"net/http"
)

var (
	UserNotFoundError     = httpErr.NewErrorResponse("USER_NOT_FOUND", http.StatusNotFound, "User not found")
	UserNameConflictError = httpErr.NewErrorResponse("USER_NAME_CONFLICT", http.StatusConflict, "User by given name already exists")
	UserOidcConflictError = httpErr.NewErrorResponse("USER_ALREADY_REGISTERED", http.StatusForbidden, "User is already registered")
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
	case error2.ErrorCodeUserNotFound:
		return UserNotFoundError
	case error2.ErrorCodeNameConflict:
		return UserNameConflictError
	case error2.ErrorCodeSubConflict:
		return UserOidcConflictError
	}
	//TODO log error
	return &err
}
