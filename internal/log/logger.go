package log

import (
	"io"
	"log/slog"
	"time"

	"idp/internal/config"

	"github.com/lmittmann/tint"
	slogformatter "github.com/samber/slog-formatter"
	slogmulti "github.com/samber/slog-multi"
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
		// TODO create lib unifying context values, slog.Logger (json, flattened), and OTEL semantics
		// inspiré de codeberg.org/shimeoki/line, davantage intégré avec slog + sémantique OTEL
		//
		// For now this shit will have to do, but it is really incomfortable
		slogHandler = slogmulti.Pipe(
			slogformatter.FlattenFormatterMiddlewareOptions{
				Separator:  ".",
				Prefix:     "",
				IgnorePath: false,
			}.NewFlattenFormatterMiddlewareOptions(),
		).Handler(
			slog.NewJSONHandler(w, &slog.HandlerOptions{
				AddSource:   cfg.Debug,
				Level:       level,
				ReplaceAttr: cfg.LoggingSchema.ReplaceAttr,
			}),
		)
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
