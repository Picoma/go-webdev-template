package server_test

// This file documents observable middleware guarantees.
//
// Middleware implementations are third-party dependencies. These tests avoid
// verifying their internal behavior and instead verify the HTTP guarantees
// relied upon by the application:
//   - configured CORS policies are applied.

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"idp/internal/config"
	"idp/internal/server"

	"github.com/stretchr/testify/require"
)

func middlewareTestRouter(t *testing.T) http.Handler {
	t.Helper()

	cfg := config.Defaults(config.Service{
		Name: "idp_test",
	})

	return server.NewRouter(
		slog.New(slog.DiscardHandler),
		cfg.WithDebug,
		&fakeSystemHandler{},
		&fakeCounterHandler{},
	)
}

func TestMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("allows configured cors origins", func(t *testing.T) {
		t.Parallel()

		router := middlewareTestRouter(t)

		req := httptest.NewRequest(http.MethodOptions, "/counter", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", http.MethodPost)
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(
			t,
			"https://example.com",
			rec.Header().Get("Access-Control-Allow-Origin"),
		)

		require.Contains(
			t,
			rec.Header().Get("Access-Control-Allow-Methods"),
			http.MethodPost,
		)

		require.Contains(
			t,
			rec.Header().Get("Access-Control-Allow-Headers"),
			"Content-Type",
		)
	})

	t.Run("rejects unsupported cors origins", func(t *testing.T) {
		t.Parallel()

		router := middlewareTestRouter(t)

		req := httptest.NewRequest(http.MethodGet, "/counter", nil)
		req.Header.Set("Origin", "ftp://example.com")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Empty(
			t,
			rec.Header().Get("Access-Control-Allow-Origin"),
		)
	})
}
