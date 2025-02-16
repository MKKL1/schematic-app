package post

import (
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/apperr"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrorCodePostNotFound apperr.Code = "POST_NOT_FOUND"
)

func ErrorMapper(err error) error {
	var code codes.Code
	if appErr, ok := apperr.FromError(err); ok {
		switch appErr.Code {
		case ErrorCodePostNotFound:
			code = codes.NotFound
		case category.ErrorCodeCategoryNotFound:
			code = codes.NotFound
		default:
			code = codes.Unknown
		}

		st := status.New(code, appErr.Message)

		errorInfo := &errdetails.ErrorInfo{
			Reason:   string(appErr.Code),
			Domain:   "schem.post",
			Metadata: appErr.Metadata,
		}
		stWithDetails, errDetails := st.WithDetails(errorInfo)
		if errDetails != nil {
			// In the rare case that attaching details fails, log and return the original status.
			return st.Err()
		}
		return stWithDetails.Err()
	}

	var pme *PostMetadataError
	if errors.As(err, &pme) {
		br := &errdetails.BadRequest{}

		for categ, v := range pme.Errors {
			for _, k := range v.Errors {
				br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
					Field:       fmt.Sprintf("%s:%s", categ, k.Field),
					Description: k.Message,
				})
			}
		}

		st := status.New(codes.InvalidArgument, "invalid post metadata")
		errorInfo := &errdetails.ErrorInfo{
			Reason:   "POST_METADATA_VALIDATION_ERROR",
			Domain:   "schem.post",
			Metadata: nil,
		}
		stWithDetails, errDetails := st.WithDetails(errorInfo, br)
		if errDetails != nil {
			return st.Err()
		}
		return stWithDetails.Err()
	}

	st := status.New(codes.Internal, "Internal Server Error")
	return st.Err()
}
