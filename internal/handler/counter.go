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

	counter, err := h.CountService.Get(ctx)
	if err != nil {
		err = httplog.SetError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httplog.SetAttrs(r.Context(),
		slog.Any("counter", counter),
	)

	if err := pages.CounterPage(counter).Render(ctx, w); err != nil {
		err = httplog.SetError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *CounterHandler) IncreaseCounter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	counter, err := h.CountService.Increment(ctx)
	if err != nil {
		err = httplog.SetError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httplog.SetAttrs(r.Context(),
		slog.Any("counter", counter),
	)

	if err := components.CounterValue(counter).Render(ctx, w); err != nil {
		err = httplog.SetError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
