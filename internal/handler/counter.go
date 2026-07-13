package handler

import (
	"context"
	"log/slog"
	"net/http"

	"idp/internal/web/templates/components"
	"idp/internal/web/templates/pages"

	"github.com/go-chi/httplog/v3"
)

type CounterService interface {
	Get(ctx context.Context) (int64, error)
	Increment(ctx context.Context) (int64, error)
}

// CounterHandler implements [server.CounterHandler].
type CounterHandler struct {
	CountService CounterService
}

func NewCounterHandler(cs CounterService) *CounterHandler {
	slog.Debug("creating CounterHandler")
	return &CounterHandler{
		CountService: cs,
	}
}

func (h *CounterHandler) DisplayCounter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	value, err := h.CountService.Get(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httplog.SetAttrs(ctx,
		slog.Int64("counter.value", value),
	)

	if err := pages.CounterPage(value).Render(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *CounterHandler) IncreaseCounter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	value, err := h.CountService.Increment(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httplog.SetAttrs(ctx,
		slog.Int64("counter.value", value),
	)

	if err := components.CounterValue(value).Render(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
