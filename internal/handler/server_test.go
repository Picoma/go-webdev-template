package handler_test

// This file documents the HTTP contract of SystemHandler.
//
// Unlike most handlers, SystemHandler intentionally has no service layer:
// system endpoints expose operational information rather than business data.
//
// These tests therefore use the real SQLite database implementation in memory.
// SQLite is the supported production database engine, so this provides a
// meaningful integration boundary rather than a mocked database.

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"idp/internal/config"
	"idp/internal/db"
	"idp/internal/handler"

	"github.com/stretchr/testify/require"
)

func newMemoryDB(t *testing.T) *db.DB {
	t.Helper()

	cfg := config.Defaults(config.Service{
		Name: "tdp_test",
	})
	cfg.DB.ConnString = ":memory:"

	database, err := db.Open(
		t.Context(),
		slog.Default(),
		&cfg.DB,
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, database.Close())
	})

	return database
}

func TestSystemHandler(t *testing.T) {
	t.Parallel()

	t.Run("returns healthy database status", func(t *testing.T) {
		t.Parallel()

		database := newMemoryDB(t)
		handler := handler.NewSystemHandler(database)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.DisplayDBHealth(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(
			t,
			"application/json",
			rec.Header().Get("Content-Type"),
		)

		var status db.DatabaseStatus
		require.NoError(
			t,
			json.NewDecoder(rec.Body).Decode(&status),
		)

		require.Equal(t, "up", status.Status)
		require.Equal(t, "healthy", status.Message)
	})

	t.Run("reports closed database as unavailable", func(t *testing.T) {
		t.Parallel()

		database := newMemoryDB(t)
		require.NoError(t, database.Close())

		handler := handler.NewSystemHandler(database)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.DisplayDBHealth(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var status db.DatabaseStatus
		require.NoError(
			t,
			json.NewDecoder(rec.Body).Decode(&status),
		)

		require.Equal(t, "down", status.Status)
		require.NotEmpty(t, status.Error)
	})
}
