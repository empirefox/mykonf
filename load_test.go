package mykonf

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigPath_Default(t *testing.T) {
	os.Unsetenv("TEST_SERVER_CONFIG")

	path := ConfigPath("TEST_")

	if path != "config.yaml" {
		t.Errorf("expected 'config.yaml', got %q", path)
	}
}

func TestConfigPath_FromEnv(t *testing.T) {
	t.Setenv("TEST_SERVER_CONFIG", "/custom/path/config.yaml")

	path := ConfigPath("TEST_")

	if path != "/custom/path/config.yaml" {
		t.Errorf("expected '/custom/path/config.yaml', got %q", path)
	}
}

func TestConfigPath_EmptyPrefix(t *testing.T) {
	t.Setenv("SERVER_CONFIG", "/no/prefix/config.yaml")

	path := ConfigPath("")

	if path != "/no/prefix/config.yaml" {
		t.Errorf("expected '/no/prefix/config.yaml', got %q", path)
	}
}

func TestLoadPath_SimpleConfig(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("name: testapp\nport: 9090\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Config struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Name != "testapp" {
		t.Errorf("expected Name='testapp', got %q", conf.Name)
	}

	if conf.Port != 9090 {
		t.Errorf("expected Port=9090, got %d", conf.Port)
	}
}

func TestLoadPath_EnvOverride(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("name: from_file\nport: 8080\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Setenv("APP_NAME", "from_env")
	t.Setenv("APP_PORT", "9999")

	type Config struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
	}

	var conf Config
	err := LoadPath("APP_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Name != "from_env" {
		t.Errorf("expected Name='from_env', got %q", conf.Name)
	}

	if conf.Port != 9999 {
		t.Errorf("expected Port=9999, got %d", conf.Port)
	}
}

func TestLoadPath_NestedConfig(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("database:\n  host: localhost\n  port: 5432\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
	type Config struct {
		Database Database `yaml:"database"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Database.Host != "localhost" {
		t.Errorf("expected Database.Host='localhost', got %q", conf.Database.Host)
	}

	if conf.Database.Port != 5432 {
		t.Errorf("expected Database.Port=5432, got %d", conf.Database.Port)
	}
}

func TestLoadPath_NestedEnvOverride(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("database:\n  host: localhost\n  port: 5432\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Setenv("APP_DATABASE_HOST", "remotehost")
	t.Setenv("APP_DATABASE_PORT", "3306")

	type Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
	type Config struct {
		Database Database `yaml:"database"`
	}

	var conf Config
	err := LoadPath("APP_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Database.Host != "remotehost" {
		t.Errorf("expected Database.Host='remotehost', got %q", conf.Database.Host)
	}

	if conf.Database.Port != 3306 {
		t.Errorf("expected Database.Port=3306, got %d", conf.Database.Port)
	}
}

func TestLoadPath_DefaultValues(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("name: testapp\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Config struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port" default:"8080"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Name != "testapp" {
		t.Errorf("expected Name='testapp', got %q", conf.Name)
	}

	if conf.Port != 8080 {
		t.Errorf("expected Port=8080 (default), got %d", conf.Port)
	}
}

func TestLoadPath_FileNotExist(t *testing.T) {
	type Config struct {
		Name string `yaml:"name" default:"default_name"`
	}

	t.Setenv("NOFILE_NAME", "from_env")

	var conf Config
	err := LoadPath("NOFILE_", "/nonexistent/config.yaml", &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Name != "from_env" {
		t.Errorf("expected Name='from_env', got %q", conf.Name)
	}
}

func TestLoadPath_FileNotExist_DefaultOnly(t *testing.T) {
	type Config struct {
		Name string `yaml:"name" default:"default_name"`
	}

	os.Unsetenv("NOFILE2_NAME")

	var conf Config
	err := LoadPath("NOFILE2_", "/nonexistent/config.yaml", &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Name != "default_name" {
		t.Errorf("expected Name='default_name', got %q", conf.Name)
	}
}

func TestLoadPath_DurationField(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("timeout: 30s\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Config struct {
		Timeout time.Duration `yaml:"timeout"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 30 * time.Second
	if conf.Timeout != expected {
		t.Errorf("expected Timeout=%v, got %v", expected, conf.Timeout)
	}
}

func TestLoadPath_DurationFromEnv(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Setenv("DUR_TIMEOUT", "5m")

	type Config struct {
		Timeout time.Duration `yaml:"timeout"`
	}

	var conf Config
	err := LoadPath("DUR_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 5 * time.Minute
	if conf.Timeout != expected {
		t.Errorf("expected Timeout=%v, got %v", expected, conf.Timeout)
	}
}

func TestLoadPath_SliceField(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("hosts:\n  - host1\n  - host2\n  - host3\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Config struct {
		Hosts []string `yaml:"hosts"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conf.Hosts) != 3 {
		t.Errorf("expected 3 hosts, got %d", len(conf.Hosts))
	}

	expected := []string{"host1", "host2", "host3"}
	for i, h := range expected {
		if i >= len(conf.Hosts) || conf.Hosts[i] != h {
			t.Errorf("expected Hosts[%d]=%q, got %q", i, h, conf.Hosts[i])
		}
	}
}

func TestLoadPath_SliceFromEnv(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Setenv("SLICE_TAGS", "tag1,tag2,tag3")

	type Config struct {
		Tags []string `yaml:"tags"`
	}

	var conf Config
	err := LoadPath("SLICE_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conf.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d: %v", len(conf.Tags), conf.Tags)
	}

	expected := []string{"tag1", "tag2", "tag3"}
	for i, tag := range expected {
		if conf.Tags[i] != tag {
			t.Errorf("expected Tags[%d]=%q, got %q", i, tag, conf.Tags[i])
		}
	}
}

func TestLoadPath_EnvExpansionInFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	t.Setenv("MY_SECRET", "supersecret")

	content := []byte("password: $MY_SECRET\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Config struct {
		Password string `yaml:"password"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Password != "supersecret" {
		t.Errorf("expected Password='supersecret', got %q", conf.Password)
	}
}

func TestLoadPath_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("invalid: yaml: content:\n  bad indentation")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Config struct {
		Name string `yaml:"name"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoadPath_BoolField(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("enabled: true\ndebug: false\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Config struct {
		Enabled bool `yaml:"enabled"`
		Debug   bool `yaml:"debug"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !conf.Enabled {
		t.Error("expected Enabled=true")
	}

	if conf.Debug {
		t.Error("expected Debug=false")
	}
}

func TestLoadPath_BoolFromEnv(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("enabled: false\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Setenv("BOOL_ENABLED", "true")

	type Config struct {
		Enabled bool `yaml:"enabled"`
	}

	var conf Config
	err := LoadPath("BOOL_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !conf.Enabled {
		t.Error("expected Enabled=true from env override")
	}
}

func TestLoad_UsesConfigPath(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "myconfig.yaml")

	content := []byte("name: loadtest\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Setenv("LOADTEST_SERVER_CONFIG", tmpFile)

	type Config struct {
		Name string `yaml:"name"`
	}

	var conf Config
	err := Load("LOADTEST_", &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Name != "loadtest" {
		t.Errorf("expected Name='loadtest', got %q", conf.Name)
	}
}

func TestLoadPath_PointerField(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := []byte("count: 42\n")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Config struct {
		Count *int `yaml:"count"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Count == nil {
		t.Fatal("expected Count to be non-nil")
	}

	if *conf.Count != 42 {
		t.Errorf("expected *Count=42, got %d", *conf.Count)
	}
}

func TestLoadPath_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(tmpFile, []byte{}, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	type Config struct {
		Name string `yaml:"name" default:"default"`
	}

	var conf Config
	err := LoadPath("TEST_", tmpFile, &conf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Name != "default" {
		t.Errorf("expected Name='default', got %q", conf.Name)
	}
}
