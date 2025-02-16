package validator

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic/ast"
)

// EnumValidator ensures a string is one of the provided options.
type EnumValidator struct {
	Values []string
}

func (e EnumValidator) Validate(value interface{}) (interface{}, error) {
	v, ok := value.(string)
	if !ok {
		return nil, errors.New("expected string")
	}
	for _, allowed := range e.Values {
		if v == allowed {
			return v, nil
		}
	}
	return nil, fmt.Errorf("value %q is not in allowed set %v", v, e.Values)
}

func NewEnumValidator(key string, node *ast.Node) (Validator, error) {
	valuesNode := node.Get("values")
	if !valuesNode.Exists() {
		return nil, fmt.Errorf("enum type for key %s missing 'values'", key)
	}
	values, err := valuesNode.Array()
	if err != nil {
		return nil, fmt.Errorf("error parsing enum values for key %s: %w", key, err)
	}

	valuesStr := make([]string, len(values))
	for i, value := range values {
		valuesStr[i] = fmt.Sprintf("%v", value)
	}

	return EnumValidator{Values: valuesStr}, nil
}
