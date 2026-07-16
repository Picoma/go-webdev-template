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
		AddSource: cfg.Verbose,
		Level:     level,
	})

	logger := slog.New(handler).With(slog.Any(
		"service", map[string]string{
			"name":        cfg.Service.Name,
			"version":     cfg.Service.Version,
			"hash_commit": cfg.Service.HashCommit,
			"env":         cfg.Service.Env,
		}))
	slog.SetDefault(logger)

	return logger
}
