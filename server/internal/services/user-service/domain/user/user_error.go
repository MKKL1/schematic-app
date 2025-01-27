package user

import appError "github.com/MKKL1/schematic-app/server/internal/pkg/error"

var (
	ErrCodeUserNotFound appError.ErrorCode = "USER_NOT_FOUND"
	ErrCodeNameConflict appError.ErrorCode = "USER_NAME_CONFLICT"
	ErrCodeSubConflict  appError.ErrorCode = "USER_SUB_CONFLICT"
)
