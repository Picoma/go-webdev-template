package log

import (
	"io"
	"log/slog"

	"idp/internal/config"
)

func New(w io.Writer, cfg *config.Config) *slog.Logger {
	level := slog.LevelInfo
	if cfg.Verbose {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:       level,
		AddSource:   cfg.Verbose,
		ReplaceAttr: cfg.LoggingSchema.ReplaceAttr,
	})

	logger := slog.New(handler).With(
		slog.Any("service", cfg.Service),
		slog.Any("config", cfg),
	)

	return logger
}
