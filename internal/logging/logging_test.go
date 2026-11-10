package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestLevelFiltering_WarnSuppressesDebugAndInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithWriter(&buf, "WARN")

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	out := buf.String()
	if strings.Contains(out, "debug message") {
		t.Error("WARN level logger should suppress DEBUG messages")
	}
	if strings.Contains(out, "info message") {
		t.Error("WARN level logger should suppress INFO messages")
	}
	if !strings.Contains(out, "warn message") {
		t.Error("WARN level logger should emit WARN messages")
	}
	if !strings.Contains(out, "error message") {
		t.Error("WARN level logger should emit ERROR messages")
	}
}

func TestLevelFiltering_DebugEmitsEverything(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithWriter(&buf, "DEBUG")

	logger.Debug("debug message")
	logger.Info("info message")

	out := buf.String()
	if !strings.Contains(out, "debug message") {
		t.Error("DEBUG level logger should emit DEBUG messages")
	}
	if !strings.Contains(out, "info message") {
		t.Error("DEBUG level logger should emit INFO messages")
	}
}

func TestParseLevel_UnknownDefaultsToInfo(t *testing.T) {
	if got := ParseLevel("NOT-A-LEVEL"); got != ParseLevel("INFO") {
		t.Errorf("ParseLevel(unknown) = %v, want INFO", got)
	}
}

// TestSecret_NeverLogsRawValue is the log-scrubbing test required by
// Task 0.3: it guards against JWTs/secrets ever appearing verbatim in
// operational logs.
func TestSecret_NeverLogsRawValue(t *testing.T) {
	const rawToken = "eyJhbGciOiJIUzI1NiJ9.super-secret-payload.signature"

	var buf bytes.Buffer
	logger := NewWithWriter(&buf, "DEBUG")

	logger.Info("token validation failed", "token", Secret(rawToken))

	out := buf.String()
	if strings.Contains(out, rawToken) {
		t.Fatalf("logged output contains the raw secret value: %s", out)
	}
	if !strings.Contains(out, redacted) {
		t.Fatalf("logged output does not contain the redaction marker: %s", out)
	}
}
