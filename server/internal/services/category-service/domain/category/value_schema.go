package category

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// RawPropertySchema holds the type information as well as any extra options.
type RawPropertySchema struct {
	Type    string                     `json:"type"`
	Options map[string]json.RawMessage `json:"-"` // any extra configuration options
}

// UnmarshalJSON is a custom unmarshaler that extracts the "type" key and stores all other keys in Options.
func (r *RawPropertySchema) UnmarshalJSON(data []byte) error {
	// Unmarshal into a temporary map.
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Extract and remove the "type" key.
	typeRaw, ok := temp["type"]
	if !ok {
		return errors.New("missing required field 'type'")
	}
	if err := json.Unmarshal(typeRaw, &r.Type); err != nil {
		return err
	}
	delete(temp, "type")

	// Store the remaining keys in Options.
	r.Options = temp

	return nil
}

type Validator interface {
	// Validate returns an error if the value does not conform.
	Validate(value interface{}) error
}

// BooleanValidator ensures the value is a boolean.
type BooleanValidator struct{}

func (b BooleanValidator) Validate(value interface{}) error {
	if _, ok := value.(bool); !ok {
		return errors.New("expected boolean")
	}
	return nil
}

// EnumValidator ensures a string is one of the provided options.
type EnumValidator struct {
	Values []string
}

func (e EnumValidator) Validate(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("expected string")
	}
	for _, allowed := range e.Values {
		if v == allowed {
			return nil
		}
	}
	return fmt.Errorf("value %q is not in allowed set %v", v, e.Values)
}

// RangeValidator validates an object with "min" and "max" numeric fields.
type RangeValidator struct {
	AllowedMin float64
	AllowedMax float64
}

func (r RangeValidator) Validate(value interface{}) error {
	// Expecting an object with "min" and "max".
	obj, ok := value.(map[string]interface{})
	if !ok {
		return errors.New("expected object with 'min' and 'max'")
	}
	vMin, ok := obj["min"]
	if !ok {
		return errors.New("missing 'min' in value")
	}
	vMax, ok := obj["max"]
	if !ok {
		return errors.New("missing 'max' in value")
	}
	minFloat, err := toFloat64(vMin)
	if err != nil {
		return fmt.Errorf("invalid 'min': %w", err)
	}
	maxFloat, err := toFloat64(vMax)
	if err != nil {
		return fmt.Errorf("invalid 'max': %w", err)
	}
	if minFloat < r.AllowedMin || maxFloat > r.AllowedMax {
		return fmt.Errorf("value range [%f, %f] not within allowed [%f, %f]",
			minFloat, maxFloat, r.AllowedMin, r.AllowedMax)
	}
	return nil
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

// ValueSchema maps a property name to its Validator.
type ValueSchema map[string]Validator

// CompileSchema compiles the raw JSON schema (from the DB) into a ValueSchema.
// It uses the "type" field and the flexible Options map to instantiate validators.
func CompileSchema(rawSchema []byte) (ValueSchema, error) {
	// Unmarshal into a map of property name to RawPropertySchema.
	var rawMap map[string]RawPropertySchema
	if err := json.Unmarshal(rawSchema, &rawMap); err != nil {
		return nil, err
	}

	compiled := make(ValueSchema)
	for key, rawProp := range rawMap {
		switch rawProp.Type {
		case "boolean":
			compiled[key] = BooleanValidator{}

		case "enum":
			// For enums, we expect an option called "values" which is an array of strings.
			valuesRaw, ok := rawProp.Options["values"]
			if !ok {
				return nil, fmt.Errorf("enum type for key %s missing 'values' option", key)
			}
			var values []string
			if err := json.Unmarshal(valuesRaw, &values); err != nil {
				return nil, fmt.Errorf("error parsing enum values for key %s: %w", key, err)
			}
			compiled[key] = EnumValidator{Values: values}

		case "range":
			// For ranges, we expect "min" and "max" options.
			minRaw, okMin := rawProp.Options["min"]
			maxRaw, okMax := rawProp.Options["max"]
			if !okMin || !okMax {
				return nil, fmt.Errorf("range type for key %s must have both 'min' and 'max' options", key)
			}
			var allowedMin, allowedMax float64
			if err := json.Unmarshal(minRaw, &allowedMin); err != nil {
				return nil, fmt.Errorf("invalid 'min' for key %s: %w", key, err)
			}
			if err := json.Unmarshal(maxRaw, &allowedMax); err != nil {
				return nil, fmt.Errorf("invalid 'max' for key %s: %w", key, err)
			}
			compiled[key] = RangeValidator{
				AllowedMin: allowedMin,
				AllowedMax: allowedMax,
			}

		default:
			return nil, fmt.Errorf("unknown type %q for key %s", rawProp.Type, key)
		}
	}
	return compiled, nil
}

// ValidateData uses the compiled schema to check the provided data.
func (vs ValueSchema) ValidateData(data map[string]interface{}) error {
	for key, validator := range vs {
		value, exists := data[key]
		if !exists {
			return fmt.Errorf("missing key: %s", key)
		}
		if err := validator.Validate(value); err != nil {
			return fmt.Errorf("validation failed for key %s: %w", key, err)
		}
	}
	return nil
}
