package server_test

// This file documents the HTTP routing contract of the application.
//
// Router tests intentionally use fake handler implementations. The goal is to
// verify that URLs are connected to the correct HTTP handlers, not to re-test
// handler behavior, templates, databases, or business logic.
//
// This creates a clean acceptance boundary:
//
//     HTTP request
//          ↓
//     middleware
//          ↓
//     router
//          ↓
//     handler interface

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"idp/internal/config"
	"idp/internal/server"

	"github.com/stretchr/testify/require"
)

func newTestRouter(t *testing.T, withDebug bool) (
	http.Handler,
	*fakeSystemHandler,
	*fakeCounterHandler,
) {
	t.Helper()

	cfg := config.Defaults(config.Service{
		Name: "idp_test",
	})
	cfg.WithDebug = withDebug

	system := &fakeSystemHandler{}
	counter := &fakeCounterHandler{}

	router := server.NewRouter(
		slog.New(slog.DiscardHandler),
		cfg,
		system,
		counter,
	)

	return router, system, counter
}

// TestRouter verifies the public HTTP routing contract.
//
// Routes are described as data because the router's responsibility is exactly
// to map HTTP requests to handlers. The fake handlers allow this test to verify
// routing without duplicating handler behavior.
func TestRouter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		method               string
		path                 string
		withDebug            bool
		wantStatus           int
		wantLocation         string
		wantCounterDisplay   bool
		wantCounterIncrement bool
		wantSystemHealth     bool
	}{
		{
			name:         "redirects root to counter",
			method:       http.MethodGet,
			path:         "/",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "/counter",
		},
		{
			name:               "routes counter display",
			method:             http.MethodGet,
			path:               "/counter",
			wantStatus:         http.StatusNoContent,
			wantCounterDisplay: true,
		},
		{
			name:                 "routes counter increment",
			method:               http.MethodPost,
			path:                 "/counter",
			wantStatus:           http.StatusNoContent,
			wantCounterIncrement: true,
		},
		{
			name:             "routes database health endpoint",
			method:           http.MethodGet,
			path:             "/health",
			wantStatus:       http.StatusNoContent,
			wantSystemHealth: true,
		},
		{
			name:       "debug routes are disabled",
			method:     http.MethodGet,
			path:       "/debug/pprof",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "debug routes are enabled",
			method:     http.MethodGet,
			path:       "/debug/pprof/",
			withDebug:  true,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router, system, counter := newTestRouter(t, tt.withDebug)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantLocation != "" {
				require.Equal(
					t,
					tt.wantLocation,
					rec.Header().Get("Location"),
				)
			}

			require.Equal(
				t,
				tt.wantCounterDisplay,
				counter.displayCalled,
			)

			require.Equal(
				t,
				tt.wantCounterIncrement,
				counter.incrementCalled,
			)

			require.Equal(
				t,
				tt.wantSystemHealth,
				system.healthCalled,
			)
		})
	}
}
