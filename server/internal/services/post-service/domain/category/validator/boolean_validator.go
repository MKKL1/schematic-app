package validator

import (
	"errors"
	"github.com/bytedance/sonic/ast"
)

// BooleanValidator ensures the value is a boolean.
type BooleanValidator struct{}

func (b BooleanValidator) Validate(value interface{}) (interface{}, error) {
	v, ok := value.(bool)
	if !ok {
		return nil, errors.New("expected boolean")
	}
	return v, nil
}

func NewBooleanValidator(key string, node *ast.Node) (Validator, error) {
	return BooleanValidator{}, nil
}
