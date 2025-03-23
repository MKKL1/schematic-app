package file

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//TODO separate generic parts

func ErrorMapper(err error) error {
	var code codes.Code
	if appErr, ok := apperr.FromError(err); ok {
		switch appErr.Code {
		case apperr.ErrorCodeBadRequest:
			code = codes.InvalidArgument
		default:
			code = codes.Unknown
		}

		st := status.New(code, appErr.Message)

		errorInfo := &errdetails.ErrorInfo{
			Reason:   string(appErr.Code),
			Domain:   "schem.file",
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
