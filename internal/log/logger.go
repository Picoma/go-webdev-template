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
	if cfg.Verbose {
		level = slog.LevelDebug
	}

	var slogHandler slog.Handler
	if cfg.TintedLogs {
		slogHandler = tint.NewHandler(w, &tint.Options{
			Level:      level,
			AddSource:  cfg.Verbose,
			TimeFormat: time.TimeOnly,
		})
	} else {
		// TODO create lib implementing wide events through OTEL spans + span events
		//
		// For now this will have to do, but it is really incomfortable for its intended purpose
		slogHandler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource: cfg.Verbose,
			Level:     level,
		})
	}

	logger := slog.New(slogHandler).With(slog.Any(
		"service", map[string]string{
			"name":        cfg.Service.Name,
			"version":     cfg.Service.Version,
			"hash_commit": cfg.Service.Commit,
			"env":         cfg.Service.Env,
		}))
	slog.SetDefault(logger)

	return logger
}
