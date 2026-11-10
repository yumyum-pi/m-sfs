// Package logging provides the application's structured operational
// logger (doc/README.md §24). It is a thin wrapper around log/slog that
// wires up the configured level and gives callers a redaction convention
// for sensitive values (JWTs, secrets) so they never end up in log output.
package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

const redacted = "***REDACTED***"

// New builds a structured JSON logger filtered to the given level
// ("DEBUG", "INFO", "WARN", or "ERROR", case-insensitive). Unrecognized
// levels fall back to INFO.
func New(level string) *slog.Logger {
	return NewWithWriter(os.Stdout, level)
}

// NewWithWriter is like New but writes to an arbitrary destination.
// Primarily useful for tests that need to inspect log output.
func NewWithWriter(w io.Writer, level string) *slog.Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: ParseLevel(level)})
	return slog.New(handler)
}

// ParseLevel converts a config-style level string into a slog.Level,
// defaulting to Info for unrecognized input.
func ParseLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Secret wraps a sensitive value (a JWT, password, API key, etc.) so that
// it is always redacted wherever it is logged, regardless of handler or
// verbosity level.
//
// Hard requirement (§24): JWT tokens and other sensitive information must
// never be written to logs. Every later phase that logs a token or secret
// must wrap it in Secret rather than passing the raw value, e.g.:
//
//	logger.Warn("token validation failed", "token", logging.Secret(rawToken))
type Secret string

// LogValue implements slog.LogValuer so the redacted form is used by every
// slog handler (JSON, text, or custom), not just a specific call site.
func (Secret) LogValue() slog.Value {
	return slog.StringValue(redacted)
}

var _ slog.LogValuer = Secret("")
