package validator

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic/ast"
	"reflect"
)

// RangeValidator validates an object with "min" and "max" numeric fields.
type RangeValidator struct {
	AllowedMin float64
	AllowedMax float64
}

func (r RangeValidator) Validate(value interface{}) (interface{}, error) {
	// Expecting an object with "min" and "max".
	obj, ok := value.(map[string]interface{})
	if !ok {
		return nil, errors.New("expected object with 'min' and 'max'")
	}
	vMin, ok := obj["min"]
	if !ok {
		return nil, errors.New("missing 'min' in value")
	}
	vMax, ok := obj["max"]
	if !ok {
		return nil, errors.New("missing 'max' in value")
	}
	minFloat, err := toFloat64(vMin)
	if err != nil {
		return nil, fmt.Errorf("invalid 'min': %w", err)
	}
	maxFloat, err := toFloat64(vMax)
	if err != nil {
		return nil, fmt.Errorf("invalid 'max': %w", err)
	}
	if minFloat < r.AllowedMin || maxFloat > r.AllowedMax {
		return nil, fmt.Errorf("value range [%f, %f] not within allowed [%f, %f]",
			minFloat, maxFloat, r.AllowedMin, r.AllowedMax)
	}
	// Return a new map with only the allowed fields.
	cleaned := map[string]interface{}{
		"min": minFloat,
		"max": maxFloat,
	}
	return cleaned, nil
}

func NewRangeValidator(key string, node *ast.Node) (Validator, error) {
	minNode := node.Get("min")
	maxNode := node.Get("max")
	if !minNode.Exists() || !maxNode.Exists() {
		return nil, fmt.Errorf("range type for key %s must have both 'min' and 'max'", key)
	}
	min, err := minNode.Float64()
	if err != nil {
		return nil, fmt.Errorf("invalid 'min' for key %s: %w", key, err)
	}
	max, err := maxNode.Float64()
	if err != nil {
		return nil, fmt.Errorf("invalid 'max' for key %s: %w", key, err)
	}
	return RangeValidator{
		AllowedMin: min,
		AllowedMax: max,
	}, nil
}

// toFloat64 converts a numeric value (which may be int or float64) to float64.
func toFloat64(v interface{}) (float64, error) {
	switch num := v.(type) {
	case float64:
		return num, nil
	case int:
		return float64(num), nil
	default:
		return 0, fmt.Errorf("got %v (type %s)", v, reflect.TypeOf(v))
	}
}
