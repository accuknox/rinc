package util

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/accuknox/rinc/internal/conf"
)

// IsosecLayout defines isosec timestamp.
const IsosecLayout = "20060102150405"

// GetNamespaceFromFQDN parses a Kubernetes FQDN and extracts the namespace
// from it.
func GetNamespaceFromFQDN(fqdn string) (string, error) {
	tokens := strings.Split(fqdn, ".")
	if len(tokens) < 2 {
		return "", fmt.Errorf("%q is not a fqdn", fqdn)
	}
	return tokens[1], nil
}

// NewLogger creates a new slog logger from the provided configuration.
func NewLogger(c conf.Log) *slog.Logger {
	var level slog.Level
	switch c.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	opt := &slog.HandlerOptions{Level: level}
	if c.Format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stderr, opt))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, opt))
}

// FileExists checks whether a file exists at the specified path.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("getting file %q info: %w", path, err)
	}
	return true, nil
}

// IsIsosec checks if the given string is an isosec identifier.
func IsIsosec(s string) bool {
	r := regexp.MustCompile("^[0-9]{14}$")
	return r.MatchString(s)
}
