package http

import (
	"errors"
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	httpErr "github.com/MKKL1/schematic-app/server/internal/pkg/http"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/services"
	"net/http"
)

var (
	UserNotFoundError     = httpErr.NewErrorResponse("USER_NOT_FOUND", http.StatusNotFound, "User not found")
	UserNameConflictError = httpErr.NewErrorResponse("USER_NAME_CONFLICT", http.StatusConflict, "User by given name already exists")
	UserOidcConflictError = httpErr.NewErrorResponse("USER_ALREADY_REGISTERED", http.StatusForbidden, "User is already registered")
)

func MapError(err error) error {
	var e *appErr.Error
	if errors.As(err, &e) {
		return MapAppError(*e)
	}
	return err
}

func MapAppError(err appErr.Error) error {
	switch err.Code() {
	case services.ErrorCodeUserNotFound:
		return UserNotFoundError
	case services.ErrorCodeNameConflict:
		return UserNameConflictError
	case services.ErrorCodeSubConflict:
		return UserOidcConflictError
	}
	//TODO log error
	return &err
}
