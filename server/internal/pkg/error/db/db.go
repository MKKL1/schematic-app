package db

import (
	"errors"
	"fmt"
)

var (
	ErrNoRows = errors.New("no rows in result set")
)

type UniqueConstraintError struct {
	Field string
}

func (e *UniqueConstraintError) Error() string {
	return fmt.Sprintf("unique constraint violated on field: %s", e.Field)
}

func NewUniqueConstraintError(field string) *UniqueConstraintError {
	return &UniqueConstraintError{Field: field}
}
