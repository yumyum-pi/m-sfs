package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeConfig(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	return path
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if err == nil {
		t.Fatal("expected an error for a missing config file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeConfig(t, `{ not valid json`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected an error for invalid JSON, got nil")
	}
}

func TestLoad_SupportedFilesOverrideReplacesDefault(t *testing.T) {
	path := writeConfig(t, `{
		"files": { "supported_files": ["png", "webp"] },
		"logging": { "level": "INFO" }
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"png", "webp"}
	if !reflect.DeepEqual(cfg.Files.SupportedFiles, want) {
		t.Fatalf("supported_files = %v, want exactly %v (no defaults merged in)", cfg.Files.SupportedFiles, want)
	}
}

func TestLoad_SupportedFilesDefaultWhenOmitted(t *testing.T) {
	path := writeConfig(t, `{ "logging": { "level": "INFO" } }`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"svg", "png", "jpeg", "webp", "gif", "pdf"}
	if !reflect.DeepEqual(cfg.Files.SupportedFiles, want) {
		t.Fatalf("supported_files = %v, want default %v", cfg.Files.SupportedFiles, want)
	}
}

func TestLoad_MaxFileSizeParsing(t *testing.T) {
	cases := []struct {
		input string
		want  int64
	}{
		{"5MB", 5 * 1024 * 1024},
		{"512KB", 512 * 1024},
		{"1GB", 1024 * 1024 * 1024},
		{"100", 100},
	}

	for _, tc := range cases {
		path := writeConfig(t, `{
			"files": { "max_file_size": "`+tc.input+`" },
			"logging": { "level": "INFO" }
		}`)

		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.input, err)
		}
		if cfg.Files.MaxFileSizeBytes != tc.want {
			t.Fatalf("input %q: MaxFileSizeBytes = %d, want %d", tc.input, cfg.Files.MaxFileSizeBytes, tc.want)
		}
	}
}

func TestLoad_MaxFileSizeDefault(t *testing.T) {
	path := writeConfig(t, `{ "logging": { "level": "INFO" } }`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Files.MaxFileSize != "5MB" {
		t.Fatalf("MaxFileSize default = %q, want %q", cfg.Files.MaxFileSize, "5MB")
	}
	if cfg.Files.MaxFileSizeBytes != 5*1024*1024 {
		t.Fatalf("MaxFileSizeBytes default = %d, want %d", cfg.Files.MaxFileSizeBytes, 5*1024*1024)
	}
}

func TestLoad_InvalidMaxFileSize(t *testing.T) {
	path := writeConfig(t, `{
		"files": { "max_file_size": "not-a-size" },
		"logging": { "level": "INFO" }
	}`)

	if _, err := Load(path); err == nil {
		t.Fatal("expected an error for an invalid max_file_size, got nil")
	}
}

func TestLoad_InvalidLoggingLevel(t *testing.T) {
	path := writeConfig(t, `{ "logging": { "level": "VERBOSE" } }`)

	if _, err := Load(path); err == nil {
		t.Fatal("expected an error for an invalid logging.level, got nil")
	}
}

func TestLoad_LoggingLevelDefault(t *testing.T) {
	path := writeConfig(t, `{}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Logging.Level != "INFO" {
		t.Fatalf("Logging.Level default = %q, want %q", cfg.Logging.Level, "INFO")
	}
}

func TestLoad_FullConfig(t *testing.T) {
	path := writeConfig(t, `{
		"files": {
			"target_dir": "/var/static/assets",
			"temp_dir": "/var/static/tmp",
			"max_file_size": "5MB",
			"supported_files": ["svg", "png", "jpeg", "webp", "gif", "pdf"]
		},
		"jwt": {
			"signing_method": "HS256",
			"secret": "test-secret",
			"issuer": "example",
			"audience": "example-api"
		},
		"database": {
			"driver": "sqlite",
			"path": "./data/dev.db"
		},
		"logging": { "level": "DEBUG" }
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Files.TargetDir != "/var/static/assets" {
		t.Errorf("TargetDir = %q", cfg.Files.TargetDir)
	}
	if cfg.Database.Driver != "sqlite" {
		t.Errorf("Database.Driver = %q", cfg.Database.Driver)
	}
	if cfg.JWT.SigningMethod != "HS256" {
		t.Errorf("JWT.SigningMethod = %q", cfg.JWT.SigningMethod)
	}
	if cfg.Logging.Level != "DEBUG" {
		t.Errorf("Logging.Level = %q", cfg.Logging.Level)
	}
}
