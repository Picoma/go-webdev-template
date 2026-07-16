package log

import (
	"context"
	"io"
	"log/slog"

	"idp/internal/config"
)

func NewWithContext(ctx context.Context, w io.Writer, cfg *config.Config) (context.Context, *slog.Logger) {
	level := slog.LevelInfo
	if cfg.Verbose {
		level = slog.LevelDebug
	}

	handler := &ContextHandler{
		next: slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource: cfg.Verbose,
			Level:     level,
		}),
	}
	logger := slog.New(handler)

	ctx = context.WithValue(ctx, ctxKeyLogAttrs{}, &[]slog.Attr{
		slog.String("service.name", cfg.Service.Name),
		slog.String("service.version", cfg.Service.Version),
		slog.String("service.hash_commit", cfg.Service.HashCommit),
		slog.String("service.env", cfg.Service.Env),
	})

	return ctx, logger
}
