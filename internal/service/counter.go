package service

import (
	"context"
	"fmt"
	"log/slog"

	"idp/internal/db/queries"
	"idp/internal/log"

	"github.com/go-chi/httplog/v3"
)

type CounterStore interface {
	GetCounter(ctx context.Context) (int64, error)
	IncrementAndGetCounter(ctx context.Context) (int64, error)
}

// CountService implements the service layer.
type CountService struct {
	store CounterStore
}

func NewCountService(q CounterStore) *CountService {
	slog.Debug("creating CountService")
	return &CountService{
		store: q,
	}
}

func (s *CountService) Get(ctx context.Context) (int64, error) {
	value, err := s.store.GetCounter(ctx)
	httplog.SetAttrs(ctx, log.SQLSchema(
		queries.GetCounter,
		"counter",
		"SELECT",
		"GetCounter",
		1,
	))
	if err != nil {
		return 0, fmt.Errorf("error getting counter: %w", err)
	}

	return value, nil
}

func (s *CountService) Increment(ctx context.Context) (int64, error) {
	value, err := s.store.IncrementAndGetCounter(ctx)
	httplog.SetAttrs(ctx, log.SQLSchema(
		queries.IncrementAndGetCounter,
		"counter",
		"INSERT",
		"IncrementAndGetCounter",
		1,
	))
	if err != nil {
		return 0, fmt.Errorf("error incrementing counter: %w", err)
	}

	return value, nil
}
