package log

import (
	"io"
	"log/slog"

	"idp/internal/config"
	"idp/internal/log/schema"
)

func New(w io.Writer, cfg *config.Config) *slog.Logger {
	level := slog.LevelInfo
	if cfg.Verbose {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource: cfg.Verbose,
		Level:     level,
	})

	logger := slog.New(handler).With(
		schema.Service(cfg.Service),
	)
	slog.SetDefault(logger)

	return logger
}
