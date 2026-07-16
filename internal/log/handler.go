package log

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	next slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := getContextAttrs(ctx)
	r.AddAttrs(attrs...)
	//nolint:wrapcheck // no errors possible from context attrs
	return h.next.Handle(ctx, r)
}

func (h *ContextHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.next.Enabled(ctx, l)
}
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{next: h.next.WithAttrs(attrs)}
}
func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{next: h.next.WithGroup(name)}
}
