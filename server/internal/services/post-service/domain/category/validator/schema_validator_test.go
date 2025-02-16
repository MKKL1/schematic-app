package validator

import (
	"github.com/bytedance/sonic"
	"os"
	"strings"
	"testing"
)

func compileTestSchema(t *testing.T, schemaStr string) MetadataSchemaValidator {
	t.Helper()
	schema, err := CompileSchema([]byte(schemaStr))
	if err != nil {
		t.Fatalf("failed to compile schema: %v", err)
	}
	return schema
}

const validSchema = `{
	"afkable": {"type": "boolean"},
	"mob_type": {"type": "enum", "values": ["zombie", "skeleton", "spider", "creeper"]},
	"spawn_rate": {"min": 100, "max": 500, "type": "range"}
}`

func TestCompileSchema(t *testing.T) {
	schema := compileTestSchema(t, validSchema)
	expectedKeys := []string{"afkable", "mob_type", "spawn_rate"}
	for _, key := range expectedKeys {
		if _, ok := schema.Map()[key]; !ok {
			t.Fatalf("key %s not found in compiled schema", key)
		}
	}
}

func TestValidateMetadata(t *testing.T) {
	schema := compileTestSchema(t, validSchema)
	data := map[string]interface{}{
		"afkable":    false,
		"mob_type":   "spider",
		"spawn_rate": map[string]interface{}{"min": 342, "max": 440, "extra": "should be removed"},
	}
	cleanedData, err := schema.ValidateData(data)
	if err != nil {
		t.Fatalf("validation error: %v", err)
	}

	// Check that only allowed keys are present.
	if len(cleanedData) != 3 {
		t.Fatalf("unexpected number of keys in validated data: got %d, want %d", len(cleanedData), 3)
	}

	// Check that the spawn_rate field is cleaned.
	spawnRate, ok := cleanedData["spawn_rate"].(map[string]interface{})
	if !ok {
		t.Fatalf("spawn_rate is not a valid object")
	}
	if _, exists := spawnRate["extra"]; exists {
		t.Fatal("extra field was not removed from spawn_rate")
	}
}

func TestInvalidSchemasFromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/invalid_schemas.json")
	if err != nil {
		t.Fatalf("failed to load test file: %v", err)
	}

	var testCases []struct {
		Name          string `json:"name"`
		Schema        string `json:"schema"`
		ExpectedError string `json:"expected_error"`
	}
	if err := sonic.Unmarshal(data, &testCases); err != nil {
		t.Fatalf("failed to unmarshal test cases: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := CompileSchema([]byte(tc.Schema))
			if err == nil {
				t.Fatalf("expected error for schema %q, got nil", tc.Schema)
			}
			if !strings.Contains(err.Error(), tc.ExpectedError) {
				t.Fatalf("expected error to contain %q, got %q", tc.ExpectedError, err.Error())
			}
		})
	}
}

func BenchmarkCompileSchema(b *testing.B) {
	_, err := CompileSchema([]byte(validSchema))
	if err != nil {
		b.Fatal(err)
	}
}
