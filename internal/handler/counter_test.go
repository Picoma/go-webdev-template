package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"idp/internal/handler"
)

// fakeCounterService implements [CounterService].
type fakeCounterService struct {
	getFunc       func(context.Context) (int64, error)
	incrementFunc func(context.Context) (int64, error)
}

func (f *fakeCounterService) Get(ctx context.Context) (int64, error) {
	return f.getFunc(ctx)
}

func (f *fakeCounterService) Increment(ctx context.Context) (int64, error) {
	return f.incrementFunc(ctx)
}

func TestNewCounterHandler(t *testing.T) {
	t.Parallel()

	service := &fakeCounterService{}

	handler := handler.NewCounterHandler(service)

	if handler.CountService != service {
		t.Fatal("expected service to be assigned")
	}
}

func TestCounterHandler_GetCounter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		value      int64
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			value:      42,
			wantStatus: http.StatusOK,
			wantBody:   "42",
		},
		{
			name:       "service error",
			err:        errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "database unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := &fakeCounterService{
				getFunc: func(context.Context) (int64, error) {
					return tt.value, tt.err
				},
			}

			handler := handler.NewCounterHandler(service)

			req := httptest.NewRequest(http.MethodGet, "/counter", nil)
			rec := httptest.NewRecorder()

			handler.DisplayCounter(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d want %d", rec.Code, tt.wantStatus)
			}

			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Fatalf("body: got %q want containing %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestCounterHandler_IncreaseCounter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		value      int64
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			value:      100,
			wantStatus: http.StatusOK,
			wantBody:   "100",
		},
		{
			name:       "service error",
			err:        errors.New("increment failed"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "increment failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := &fakeCounterService{
				incrementFunc: func(context.Context) (int64, error) {
					return tt.value, tt.err
				},
			}

			handler := handler.NewCounterHandler(service)

			req := httptest.NewRequest(http.MethodPost, "/counter/increase", nil)
			rec := httptest.NewRecorder()

			handler.IncreaseCounter(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d want %d", rec.Code, tt.wantStatus)
			}

			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Fatalf("body: got %q want containing %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
