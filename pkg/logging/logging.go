package logging

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger creates a structured slog logger configured by level.
func NewLogger(level string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(level)}
	h := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(h)
}

func parseLevel(v string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(v)) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
