package validator

import (
	"fmt"
	"strings"
)

type FieldError struct {
	// Field name of field in category metadata schema
	Field string
	// Message error message
	Message string
}

type ValidationError struct {
	Errors []FieldError
}

// Error implements the error interface.
func (ve *ValidationError) Error() string {
	var msgs []string
	for _, fe := range ve.Errors {
		msgs = append(msgs, fmt.Sprintf("%s: %s", fe.Field, fe.Message))
	}
	return strings.Join(msgs, "; ")
}
