package service

import (
	"context"
	"fmt"
	"log/slog"

	"idp/internal/db/queries"

	slogchi "github.com/samber/slog-chi"
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

	// TODO wide events
	s.logSQLQuery(ctx,
		queries.GetCounter,
		"SELECT",
		"GetCounter",
		"1",
	)

	if err != nil {
		return 0, fmt.Errorf("error getting counter: %w", err)
	}

	return value, nil
}

func (s *CountService) Increment(ctx context.Context) (int64, error) {
	value, err := s.store.IncrementAndGetCounter(ctx)

	s.logSQLQuery(ctx,
		queries.IncrementAndGetCounter,
		"INSERT",
		"IncrementAndGetCounter",
		"1",
	)

	if err != nil {
		return 0, fmt.Errorf("error incrementing counter: %w", err)
	}

	return value, nil
}

// logSQLQuery is helper
// TODO replace WE.
func (*CountService) logSQLQuery(
	ctx context.Context,
	queryText string,
	opName string,
	procName string,
	returnedRows string,
) {
	slogchi.AddContextAttributes(ctx, slog.Any(
		"db", map[string]any{
			"query": map[string]string{
				"text": queryText,
			},
			"namespace": "counter",
			"operation": map[string]string{
				"name": opName,
			},
			"stored_procedure": map[string]string{
				"name": procName,
			},
			"response": map[string]string{
				"returned_rows": returnedRows,
			},
		}),
	)
}
