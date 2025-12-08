package mykonf

import (
	"testing"
	"time"
)

func TestEnvToKey_SimpleStruct(t *testing.T) {
	type Simple struct {
		Name  string `yaml:"name"`
		Value int    `yaml:"value"`
	}

	result := EnvToKey((*Simple)(nil), "yaml")

	expected := map[string]string{
		"NAME":  "name",
		"VALUE": "value",
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(result))
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("expected result[%q] = %q, got %q", k, v, result[k])
		}
	}
}

func TestEnvToKey_NestedStruct(t *testing.T) {
	type Inner struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
	type Outer struct {
		Server Inner `yaml:"server"`
	}

	result := EnvToKey((*Outer)(nil), "yaml")

	expected := map[string]string{
		"SERVER":      "server",
		"SERVER_HOST": "server.host",
		"SERVER_PORT": "server.port",
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d keys, got %d: %v", len(expected), len(result), result)
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("expected result[%q] = %q, got %q", k, v, result[k])
		}
	}
}

func TestEnvToKey_PointerField(t *testing.T) {
	type Inner struct {
		Name string `yaml:"name"`
	}
	type Outer struct {
		Inner *Inner `yaml:"inner"`
	}

	result := EnvToKey((*Outer)(nil), "yaml")

	if result["INNER_NAME"] != "inner.name" {
		t.Errorf("expected result[INNER_NAME] = 'inner.name', got %q", result["INNER_NAME"])
	}
}

func TestEnvToKey_IgnoredField(t *testing.T) {
	type Config struct {
		Name    string `yaml:"name"`
		Ignored string `yaml:"-"`
	}

	result := EnvToKey((*Config)(nil), "yaml")

	if _, ok := result["IGNORED"]; ok {
		t.Error("ignored field should not be in result")
	}

	if result["NAME"] != "name" {
		t.Errorf("expected result[NAME] = 'name', got %q", result["NAME"])
	}
}

func TestEnvToKey_NoTag(t *testing.T) {
	type Config struct {
		FieldName string
	}

	result := EnvToKey((*Config)(nil), "yaml")

	if result["FIELDNAME"] != "FieldName" {
		t.Errorf("expected result[FIELDNAME] = 'FieldName', got %q", result["FIELDNAME"])
	}
}

func TestEnvToKey_TagWithOptions(t *testing.T) {
	type Config struct {
		Name string `yaml:"name,omitempty"`
	}

	result := EnvToKey((*Config)(nil), "yaml")

	if result["NAME"] != "name" {
		t.Errorf("expected result[NAME] = 'name', got %q", result["NAME"])
	}
}

func TestEnvToKey_TimeField(t *testing.T) {
	type Config struct {
		CreatedAt time.Time     `yaml:"created_at"`
		Timeout   time.Duration `yaml:"timeout"`
	}

	result := EnvToKey((*Config)(nil), "yaml")

	if result["CREATED_AT"] != "created_at" {
		t.Errorf("expected result[CREATED_AT] = 'created_at', got %q", result["CREATED_AT"])
	}

	if result["TIMEOUT"] != "timeout" {
		t.Errorf("expected result[TIMEOUT] = 'timeout', got %q", result["TIMEOUT"])
	}
}

func TestEnvToKey_UnexportedField(t *testing.T) {
	type Config struct {
		Name     string `yaml:"name"`
		internal string `yaml:"internal"`
	}

	result := EnvToKey((*Config)(nil), "yaml")

	if _, ok := result["INTERNAL"]; ok {
		t.Error("unexported field should not be in result")
	}

	if result["NAME"] != "name" {
		t.Errorf("expected result[NAME] = 'name', got %q", result["NAME"])
	}
}

func TestEnvToKey_DeeplyNested(t *testing.T) {
	type Level3 struct {
		Value string `yaml:"value"`
	}
	type Level2 struct {
		Level3 Level3 `yaml:"level3"`
	}
	type Level1 struct {
		Level2 Level2 `yaml:"level2"`
	}

	result := EnvToKey((*Level1)(nil), "yaml")

	expected := map[string]string{
		"LEVEL2":                "level2",
		"LEVEL2_LEVEL3":         "level2.level3",
		"LEVEL2_LEVEL3_VALUE":   "level2.level3.value",
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("expected result[%q] = %q, got %q", k, v, result[k])
		}
	}
}

func TestEnvToKey_NonStructInput(t *testing.T) {
	result := EnvToKey((*string)(nil), "yaml")

	if len(result) != 0 {
		t.Errorf("expected empty result for non-struct, got %v", result)
	}
}

func TestEnvToKey_DoublePointer(t *testing.T) {
	type Config struct {
		Name string `yaml:"name"`
	}

	var nilPtr *Config
	result := EnvToKey(&nilPtr, "yaml")

	if result["NAME"] != "name" {
		t.Errorf("expected result[NAME] = 'name', got %q", result["NAME"])
	}
}

func TestEnvToKey_DifferentTag(t *testing.T) {
	type Config struct {
		Name string `json:"json_name" yaml:"yaml_name"`
	}

	resultYaml := EnvToKey((*Config)(nil), "yaml")
	resultJson := EnvToKey((*Config)(nil), "json")

	if resultYaml["YAML_NAME"] != "yaml_name" {
		t.Errorf("expected yaml result[YAML_NAME] = 'yaml_name', got %q", resultYaml["YAML_NAME"])
	}

	if resultJson["JSON_NAME"] != "json_name" {
		t.Errorf("expected json result[JSON_NAME] = 'json_name', got %q", resultJson["JSON_NAME"])
	}
}

func TestEnvToKey_EmptyStruct(t *testing.T) {
	type Empty struct{}

	result := EnvToKey((*Empty)(nil), "yaml")

	if len(result) != 0 {
		t.Errorf("expected empty result for empty struct, got %v", result)
	}
}

func TestEnvToKey_PointerToPointerField(t *testing.T) {
	type Inner struct {
		Name string `yaml:"name"`
	}
	type Outer struct {
		Inner **Inner `yaml:"inner"`
	}

	result := EnvToKey((*Outer)(nil), "yaml")

	if result["INNER_NAME"] != "inner.name" {
		t.Errorf("expected result[INNER_NAME] = 'inner.name', got %q", result["INNER_NAME"])
	}
}
