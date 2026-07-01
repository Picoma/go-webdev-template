package log

import (
	"io"
	"log/slog"
	"time"

	"idp/internal/config"

	"github.com/lmittmann/tint"
)

func New(w io.Writer, cfg *config.Config) *slog.Logger {
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}

	var slogHandler slog.Handler
	if cfg.TintedLogs {
		slogHandler = tint.NewHandler(w, &tint.Options{
			Level:      level,
			AddSource:  cfg.Debug,
			TimeFormat: time.TimeOnly,
		})
	} else {
		slogHandler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource: cfg.Debug,
			Level:     level,
			// ReplaceAttr: format.ReplaceAttr,
		})
	}

	logger := slog.New(slogHandler).With(
		slog.String("service.name", cfg.Service.Name),
		slog.String("service.version", cfg.Service.Version),
		slog.String("service.hash_commit", cfg.Service.Commit),
		slog.String("service.env", cfg.Service.Env),
	)

	slog.SetDefault(logger)

	return logger
}
