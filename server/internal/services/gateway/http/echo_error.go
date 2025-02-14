package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func MapUnmarshalError(err error) error {
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		msg := unmarshalErr.Error()
		const prefix = "struct field "
		fieldName := ""
		if idx := strings.Index(msg, prefix); idx != -1 {
			subStr := msg[idx+len(prefix):] // e.g. "PostCreateRequest.author of type int64"
			if dotIdx := strings.Index(subStr, "."); dotIdx != -1 {
				subStr = subStr[dotIdx+1:] // e.g. "author of type int64"
				if spaceIdx := strings.Index(subStr, " "); spaceIdx != -1 {
					fieldName = subStr[:spaceIdx]
				} else {
					fieldName = subStr
				}
			}
		}
		if fieldName == "" {
			fieldName = "unknown"
		}
		message := fmt.Sprintf("Field '%s' has an invalid type", fieldName)
		metadata := map[string]string{
			"parameter": fieldName,
			"details":   "invalid type",
			"code":      "invalid_type",
			"value":     unmarshalErr.Value,
		}
		return &GatewayError{
			HttpCode: 400,
			ErrResponse: ErrorResponse{
				Errors: []ErrorDetail{{
					Domain:   "gateway",
					Reason:   "VALIDATION_ERROR",
					Message:  message,
					Metadata: metadata,
				}},
			},
		}
	}

	// If the error isn't an unmarshal type error, return it as-is.
	return err
}

func EchoErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			var syntaxErr *json.SyntaxError
			var typeErr *json.UnmarshalTypeError
			var gatewayErr *GatewayError

			if errors.As(err, &syntaxErr) {
				return c.JSON(http.StatusBadRequest, ErrorResponse{
					Errors: []ErrorDetail{
						{
							Domain:   "gateway",
							Reason:   "INVALID_JSON",
							Message:  "Invalid JSON format",
							Metadata: map[string]string{"position": fmt.Sprintf("%d", syntaxErr.Offset)},
						},
					},
				})
			}

			if errors.As(err, &typeErr) {
				return c.JSON(http.StatusBadRequest, MapUnmarshalError(typeErr).(*GatewayError).ErrResponse)
			}

			if errors.As(err, &gatewayErr) {
				return c.JSON(gatewayErr.HttpCode, gatewayErr.ErrResponse)
			}

			c.Logger().Error(err)

			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Errors: []ErrorDetail{
					{
						Domain:   "gateway",
						Reason:   "INTERNAL_ERROR",
						Message:  "An unexpected error occurred",
						Metadata: map[string]string{},
					},
				},
			})
		}
		return nil
	}
}
