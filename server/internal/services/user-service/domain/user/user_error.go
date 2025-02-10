package user

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrorCodeUserNotFound apperr.Code = "USER_NOT_FOUND"
	ErrorCodeNameConflict apperr.Code = "USER_NAME_CONFLICT"
	ErrorCodeSubConflict  apperr.Code = "USER_SUB_CONFLICT"
)

func ErrorMapper(err error) error {
	var code codes.Code
	if appErr, ok := apperr.FromError(err); ok {
		switch appErr.Code {
		case ErrorCodeUserNotFound:
			code = codes.NotFound
		case ErrorCodeNameConflict:
			code = codes.AlreadyExists
		case ErrorCodeSubConflict:
			code = codes.AlreadyExists
		default:
			code = codes.Unknown
		}

		st := status.New(code, appErr.Message)

		errorInfo := &errdetails.ErrorInfo{
			Reason:   string(appErr.Code),
			Domain:   "schem.user",
			Metadata: appErr.Metadata,
		}
		stWithDetails, errDetails := st.WithDetails(errorInfo)
		if errDetails != nil {
			// In the rare case that attaching details fails, log and return the original status.
			return st.Err()
		}
		return stWithDetails.Err()
	}

	st := status.New(codes.Internal, "Internal Server Error")
	return st.Err()
}
