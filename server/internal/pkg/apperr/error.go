package apperr

import "errors"

// Code represents a domain error code.
type Code string

const (
	ErrorCodeUnknown      Code = "UNKNOWN"
	ErrorCodeConflict     Code = "CONFLICT"
	ErrorCodeBadRequest   Code = "BAD_REQUEST"
	ErrorCodeNotFound     Code = "NOT_FOUND"
	ErrorCodeUnauthorized Code = "UNAUTHORIZED"
)

// SlugError is a custom error that wraps a domain error.
type SlugError struct {
	Code     Code
	Slug     string
	Err      error
	Metadata map[string]string //Is this good idea?
}

func NewSlugError(err error, slug string) *SlugError {
	return &SlugError{
		Code: ErrorCodeUnknown,
		Slug: slug,
		Err:  err,
	}
}

func NewSlugErrorWithCode(err error, slug string, code Code) *SlugError {
	return &SlugError{
		Code: code,
		Slug: slug,
		Err:  err,
	}
}

// Error implements the error interface.
func (e *SlugError) Error() string {
	return e.Err.Error()
}

func FromError(err error) (*SlugError, bool) {
	var appErr *SlugError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func (e *SlugError) AddMetadata(key string, value string) *SlugError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value

	return e
}
