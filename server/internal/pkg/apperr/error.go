package apperr

import (
	"errors"
	"fmt"
)

// Code represents a domain error code.
type Code string

const (
	ErrorCodeUnknown    Code = "UNKNOWN"
	ErrorCodeConflict   Code = "CONFLICT"
	ErrorCodeBadRequest Code = "BAD_REQUEST"
)

// AppError is a custom error that wraps a domain error.
type AppError struct {
	Code     Code
	Message  string
	Err      error
	Metadata map[string]string
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap allows errors.Is and errors.As to work.
func (e *AppError) Unwrap() error {
	return e.Err
}

// WrapErrorf creates a new AppError wrapping the given error.
func WrapErrorf(err error, code Code, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Err:     err,
	}
}

// FromError extracts an AppError from err if possible.
func FromError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func (e *AppError) AddMetadata(key string, value string) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value

	return e
}

/*
type Code string

type AppError struct {
	Code     Code
	Message  string
	Metadata map[string]string
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError is a constructor helper.
func NewAppError(code Code, message string) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewAppErrorWithMetadata(code Code, message string, metadata map[string]string) error {}
*/
