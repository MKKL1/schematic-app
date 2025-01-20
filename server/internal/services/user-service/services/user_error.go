package services

import appError "github.com/MKKL1/schematic-app/server/internal/pkg/error"

var (
	ErrorCodeUserNotFound appError.ErrorCode = "USER_NOT_FOUND"
	ErrorCodeNameConflict appError.ErrorCode = "USER_NAME_CONFLICT"
	ErrorCodeSubConflict  appError.ErrorCode = "USER_SUB_CONFLICT"
)
