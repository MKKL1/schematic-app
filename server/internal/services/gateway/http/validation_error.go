package http

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type GatewayError struct {
	HttpCode    int
	ErrResponse ErrorResponse
}

func NewSimpleGatewayError(httpCode int, reason string, message string, metadata map[string]string) *GatewayError {
	return &GatewayError{
		HttpCode: httpCode,
		ErrResponse: ErrorResponse{
			Errors: []ErrorDetail{
				{
					Reason:   reason,
					Message:  message,
					Metadata: metadata,
				},
			},
		},
	}
}

func (e *GatewayError) Error() string {
	if len(e.ErrResponse.Errors) > 0 {
		return e.ErrResponse.Errors[0].Message
	}
	return "gateway error"
}

// getJSONFieldName uses reflection on requestData to fetch the JSON tag for a given struct field.
func getJSONFieldName(requestData interface{}, fieldName string) string {
	t := reflect.TypeOf(requestData)
	// If requestData is a pointer, get the element.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Look up the field by its struct field name.
	if f, ok := t.FieldByName(fieldName); ok {
		// Get the json tag (e.g. "name,omitempty")
		tag := f.Tag.Get("json")
		if tag != "" {
			// The tag may contain options separated by commas; use the first part.
			parts := strings.Split(tag, ",")
			if parts[0] != "-" && parts[0] != "" {
				return parts[0]
			}
		}
	}
	return strings.ToLower(fieldName)
}

func MapValidationError(validationErr error, requestData interface{}) error {
	var vErrs validator.ValidationErrors
	if errors.As(validationErr, &vErrs) {
		detailMapping := map[string]string{
			"max":      "too long",
			"min":      "too short",
			"required": "is required",
		}

		var errorDetails []ErrorDetail

		// Iterate over all field errors.
		for _, fe := range vErrs {

			jsonFieldName := getJSONFieldName(requestData, fe.StructField())

			detail, found := detailMapping[fe.Tag()]
			if !found {
				detail = fmt.Sprintf("failed '%s' validation", fe.Tag())
			}

			vBuilder := ValidationErrorBuilder{
				Parameter: jsonFieldName,
				Detail:    detail,
				Code:      fe.Tag(),
				Value:     fmt.Sprintf("%v", fe.Value()),
				Message:   fmt.Sprintf("Field '%s' %s", jsonFieldName, detail),
			}

			errorDetails = append(errorDetails, vBuilder.Build())
		}

		return &GatewayError{
			HttpCode: 400,
			ErrResponse: ErrorResponse{
				Errors: errorDetails,
			},
		}
	}

	return &GatewayError{
		HttpCode: 500,
		ErrResponse: ErrorResponse{
			Errors: []ErrorDetail{{
				Reason:   "VALIDATION_ERROR",
				Message:  validationErr.Error(),
				Metadata: map[string]string{},
			}},
		},
	}
}

type ValidationErrorBuilder struct {
	Parameter string
	Detail    string
	Code      string
	Value     string
	Message   string
}

func (v ValidationErrorBuilder) Build() ErrorDetail {
	metadata := map[string]string{}

	if v.Parameter != "" {
		metadata["parameter"] = v.Parameter
	}

	if v.Detail != "" {
		metadata["detail"] = v.Detail
	}

	if v.Code != "" {
		metadata["code"] = v.Code
	}

	if v.Value != "" {
		metadata["value"] = v.Value
	}

	return ErrorDetail{
		Reason:   "VALIDATION_ERROR",
		Message:  v.Message,
		Metadata: metadata,
	}
}
