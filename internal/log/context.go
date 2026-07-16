package log

import (
	"context"
	"log/slog"
)

type ctxKeyLogAttrs struct{}

func (c *ctxKeyLogAttrs) String() string {
	return "httplog attrs context"
}

// AddContextAttrs sets the attributes on the context.
func AddContextAttrs(ctx context.Context, attrs ...slog.Attr) {
	if ptr, ok := ctx.Value(ctxKeyLogAttrs{}).(*[]slog.Attr); ok && ptr != nil {
		*ptr = append(*ptr, attrs...)
	}
}

func getContextAttrs(ctx context.Context) []slog.Attr {
	if ptr, ok := ctx.Value(ctxKeyLogAttrs{}).(*[]slog.Attr); ok && ptr != nil {
		return *ptr
	}

	return nil
}
