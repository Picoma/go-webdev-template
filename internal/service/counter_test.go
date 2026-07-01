package service_test

import (
	"context"
	"errors"
	"testing"

	"idp/internal/service"
)

// fakeQueries implements [CounterStore].
type fakeQueries struct {
	getFn       func(context.Context) (int64, error)
	incrementFn func(context.Context) (int64, error)
}

func (f fakeQueries) GetCounter(ctx context.Context) (int64, error) {
	return f.getFn(ctx)
}

func (f fakeQueries) IncrementAndGetCounter(ctx context.Context) (int64, error) {
	return f.incrementFn(ctx)
}

func TestCountService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		queries fakeQueries
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			queries: fakeQueries{
				getFn: func(context.Context) (int64, error) {
					return 42, nil
				},
			},
			want: 42,
		},
		{
			name: "query error",
			queries: fakeQueries{
				getFn: func(context.Context) (int64, error) {
					return 0, errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := service.NewCountService(tt.queries)

			got, err := svc.Get(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCountService_Increment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		queries fakeQueries
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			queries: fakeQueries{
				incrementFn: func(context.Context) (int64, error) {
					return 43, nil
				},
			},
			want: 43,
		},
		{
			name: "query error",
			queries: fakeQueries{
				incrementFn: func(context.Context) (int64, error) {
					return 0, errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := service.NewCountService(tt.queries)

			got, err := svc.Increment(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}
