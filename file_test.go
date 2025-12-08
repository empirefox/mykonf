package mykonf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProvider(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")

	content := []byte("key: value\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	provider := Provider(tmpFile)

	if provider.File == nil {
		t.Fatal("expected non-nil File")
	}
}

func TestFile_ReadBytes_Simple(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")

	content := []byte("name: test\nvalue: 123\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	provider := Provider(tmpFile)
	result, err := provider.ReadBytes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(result) != string(content) {
		t.Errorf("expected %q, got %q", string(content), string(result))
	}
}

func TestFile_ReadBytes_EnvExpansion(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")

	t.Setenv("TEST_VAR", "expanded_value")

	content := []byte("name: $TEST_VAR\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	provider := Provider(tmpFile)
	result, err := provider.ReadBytes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "name: expanded_value\n"
	if string(result) != expected {
		t.Errorf("expected %q, got %q", expected, string(result))
	}
}

func TestFile_ReadBytes_EnvExpansionBraces(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")

	t.Setenv("MY_HOST", "localhost")
	t.Setenv("MY_PORT", "8080")

	content := []byte("url: ${MY_HOST}:${MY_PORT}\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	provider := Provider(tmpFile)
	result, err := provider.ReadBytes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "url: localhost:8080\n"
	if string(result) != expected {
		t.Errorf("expected %q, got %q", expected, string(result))
	}
}

func TestFile_ReadBytes_UnsetEnv(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")

	os.Unsetenv("UNSET_VAR")

	content := []byte("name: $UNSET_VAR\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	provider := Provider(tmpFile)
	result, err := provider.ReadBytes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "name: \n"
	if string(result) != expected {
		t.Errorf("expected %q, got %q", expected, string(result))
	}
}

func TestFile_ReadBytes_MixedContent(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")

	t.Setenv("DB_HOST", "dbserver")

	content := []byte("static: value\ndynamic: $DB_HOST\nnumber: 42\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	provider := Provider(tmpFile)
	result, err := provider.ReadBytes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "static: value\ndynamic: dbserver\nnumber: 42\n"
	if string(result) != expected {
		t.Errorf("expected %q, got %q", expected, string(result))
	}
}

func TestFile_ReadBytes_FileNotFound(t *testing.T) {
	provider := Provider("/nonexistent/path/file.yaml")
	_, err := provider.ReadBytes()

	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestFile_ReadBytes_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty.yaml")

	if err := os.WriteFile(tmpFile, []byte{}, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	provider := Provider(tmpFile)
	result, err := provider.ReadBytes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result, got %q", string(result))
	}
}

func TestFile_ReadBytes_MultipleEnvVars(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")

	t.Setenv("VAR1", "first")
	t.Setenv("VAR2", "second")
	t.Setenv("VAR3", "third")

	content := []byte("a: $VAR1\nb: $VAR2\nc: $VAR3\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	provider := Provider(tmpFile)
	result, err := provider.ReadBytes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "a: first\nb: second\nc: third\n"
	if string(result) != expected {
		t.Errorf("expected %q, got %q", expected, string(result))
	}
}
