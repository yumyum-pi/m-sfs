package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sizePattern = regexp.MustCompile(`(?i)^\s*([0-9]*\.?[0-9]+)\s*(B|KB|MB|GB|TB)?\s*$`)

var sizeUnitMultiplier = map[string]int64{
	"":   1,
	"B":  1,
	"KB": 1024,
	"MB": 1024 * 1024,
	"GB": 1024 * 1024 * 1024,
	"TB": 1024 * 1024 * 1024 * 1024,
}

// parseSize converts a human-readable size string such as "5MB" or "512KB"
// into a byte count. Units are binary (1KB = 1024B). A bare number is
// interpreted as bytes.
func parseSize(s string) (int64, error) {
	m := sizePattern.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("invalid size %q", s)
	}

	value, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size %q: %w", s, err)
	}

	unit := strings.ToUpper(m[2])
	multiplier, ok := sizeUnitMultiplier[unit]
	if !ok {
		return 0, fmt.Errorf("invalid size %q: unknown unit %q", s, m[2])
	}

	return int64(value * float64(multiplier)), nil
}
