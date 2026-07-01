package server_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"idp/internal/config"
	"idp/internal/db"
	"idp/internal/handler"
	"idp/internal/server"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()

	cfg := &config.Config{
		DB: config.DB{
			Driver:     "sqlite3",
			ConnString: ":memory:",
		},
	}

	database, err := db.Open(
		context.Background(),
		slog.Default(),
		cfg,
	)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	t.Cleanup(func() {
		_ = database.Close()
	})

	counter := handler.NewCounterHandler(
		&fakeCounterService{},
	)

	return server.NewRouter(
		slog.Default(),
		cfg,
		database,
		counter,
	)
}

type fakeCounterService struct{}

func (fakeCounterService) Get(_ context.Context) (int64, error) {
	return 1, nil
}

func (fakeCounterService) Increment(_ context.Context) (int64, error) {
	return 2, nil
}

func TestNewRouterRoutes(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	tests := []struct {
		name   string
		method string
		path   string
		status int
	}{
		{
			name:   "hello world",
			method: http.MethodGet,
			path:   "/",
			status: http.StatusMovedPermanently,
		},
		{
			name:   "health",
			method: http.MethodGet,
			path:   "/health",
			status: http.StatusOK,
		},
		{
			name:   "counter get",
			method: http.MethodGet,
			path:   "/counter",
			status: http.StatusOK,
		},
		{
			name:   "counter post",
			method: http.MethodPost,
			path:   "/counter",
			status: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(
				tt.method,
				tt.path,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.status {
				t.Fatalf(
					"status: got %d want %d body=%s",
					rec.Code,
					tt.status,
					rec.Body.String(),
				)
			}
		})
	}
}

func TestRouterCORS(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	req := httptest.NewRequest(
		http.MethodOptions,
		"/",
		nil,
	)

	req.Header.Set(
		"Origin",
		"https://example.com",
	)

	req.Header.Set(
		"Access-Control-Request-Method",
		"GET",
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Fatalf(
			"cors origin: got %q",
			got,
		)
	}
}

func TestRouterNotFound(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/does-not-exist",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"status: got %d want %d",
			rec.Code,
			http.StatusNotFound,
		)
	}

	if strings.TrimSpace(rec.Body.String()) != "404 page not found" {
		t.Fatalf(
			"body: %q",
			rec.Body.String(),
		)
	}
}
