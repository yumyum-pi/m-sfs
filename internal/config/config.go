// Package config loads the application's config.json once at startup.
// Per doc/README.md §23, there is no hot reload: the process must be
// restarted for configuration changes to take effect.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// defaultSupportedFiles is used when files.supported_files is omitted from
// config.json. Per §8, when the field IS present it replaces this list
// entirely rather than extending it.
var defaultSupportedFiles = []string{"svg", "png", "jpeg", "webp", "gif", "pdf"}

const defaultMaxFileSize = "5MB"

var validLogLevels = map[string]bool{
	"DEBUG": true,
	"INFO":  true,
	"WARN":  true,
	"ERROR": true,
}

// Config is the root application configuration, loaded from config.json.
type Config struct {
	Files    FilesConfig    `json:"files"`
	JWT      JWTConfig      `json:"jwt"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
}

// FilesConfig covers asset storage locations and upload constraints (§6, §8, §9, §13).
type FilesConfig struct {
	TargetDir      string   `json:"target_dir"`
	TempDir        string   `json:"temp_dir"`
	MaxFileSize    string   `json:"max_file_size"`
	SupportedFiles []string `json:"supported_files"`

	// MaxFileSizeBytes is derived from MaxFileSize during Load and is not
	// itself part of the JSON schema.
	MaxFileSizeBytes int64 `json:"-"`
}

// JWTConfig holds settings needed to validate incoming tokens (§17).
// Exact fields are finalized alongside Phase 2's auth middleware.
type JWTConfig struct {
	SigningMethod string `json:"signing_method"` // e.g. "HS256", "RS256"
	Secret        string `json:"secret,omitempty"`
	JWKSURL       string `json:"jwks_url,omitempty"`
	Issuer        string `json:"issuer,omitempty"`
	Audience      string `json:"audience,omitempty"`
}

// DatabaseConfig holds connection settings for either MySQL (prod) or
// SQLite (dev), selected via Driver (§18).
type DatabaseConfig struct {
	Driver string `json:"driver"` // "mysql" or "sqlite"

	// SQLite
	Path string `json:"path,omitempty"`

	// MySQL
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`

	// DSN, if set, overrides the discrete fields above for drivers that
	// accept a raw connection string.
	DSN string `json:"dsn,omitempty"`

	MaxOpenConns int `json:"max_open_conns,omitempty"`
	MaxIdleConns int `json:"max_idle_conns,omitempty"`
}

// LoggingConfig controls operational log verbosity (§24).
type LoggingConfig struct {
	Level string `json:"level"`
}

// Load reads and parses config.json from path, applies defaults, and
// validates the result. It is intended to be called once at startup.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: could not read %q: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: could not parse %q: %w", path, err)
	}

	cfg.applyDefaults()

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if len(c.Files.SupportedFiles) == 0 {
		c.Files.SupportedFiles = append([]string(nil), defaultSupportedFiles...)
	}
	if c.Files.MaxFileSize == "" {
		c.Files.MaxFileSize = defaultMaxFileSize
	}
	if c.Logging.Level == "" {
		c.Logging.Level = "INFO"
	}
}

func (c *Config) validate() error {
	size, err := parseSize(c.Files.MaxFileSize)
	if err != nil {
		return fmt.Errorf("files.max_file_size: %w", err)
	}
	c.Files.MaxFileSizeBytes = size

	level := strings.ToUpper(c.Logging.Level)
	if !validLogLevels[level] {
		return fmt.Errorf("logging.level: invalid value %q (must be DEBUG, INFO, WARN, or ERROR)", c.Logging.Level)
	}
	c.Logging.Level = level

	return nil
}
