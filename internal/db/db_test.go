//nolint:paralleltest // If swapping databases is required/desired (like for Postgres), infra reuse might be needed, thus sequential tests
package db_test

// This file documents the lifecycle contract of the database package.
//
// The tests intentionally focus on externally observable behavior:
//   - opening a database must produce a usable connection and initialized query layer;
//   - invalid database configuration must fail early;
//   - closing a database must make the connection unusable.
//
// These tests avoid checking implementation details (such as sql.Open internals or
// logger output) because those details are not part of the package contract.

import (
	"database/sql"
	"errors"
	"log/slog"
	"path/filepath"
	"testing"

	"idp/internal/config"
	"idp/internal/db"

	"github.com/stretchr/testify/require"
)

func testConfig(t *testing.T) *config.Config {
	t.Helper()

	cfg := config.Defaults(config.Service{
		Name: "test_idp",
	})
	cfg.DB.Driver = "sqlite3"
	cfg.DB.ConnString = filepath.Join(t.TempDir(), "test.db")

	return cfg
}

func discardLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

// TestOpen verifies that database initialization performs all required startup
// steps:
//
//   - the configured SQL driver can be opened;
//   - the connection can actually be reached (Ping succeeds);
//   - the generated query layer is initialized.
//
// This protects against regressions where a database handle is returned but is
// not actually usable by the application.
func TestOpen(t *testing.T) {
	t.Run("opens a valid sqlite database", func(t *testing.T) {
		ctx := t.Context()

		db, err := db.Open(ctx, discardLogger(), &testConfig(t).DB)
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, db.Close())
		})

		require.NotNil(t, db.Connection)
		require.NotNil(t, db.Queries)
		require.NoError(t, db.Connection.PingContext(ctx))
	})

	// sql.Open only validates the driver name. This regression test ensures
	// configuration errors are reported clearly instead of producing a database
	// object that can never be used.
	t.Run("rejects an unknown database driver", func(t *testing.T) {
		cfg := testConfig(t)
		cfg.DB.Driver = "does-not-exist"

		_, err := db.Open(t.Context(), discardLogger(), &cfg.DB)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to open database")
	})

	// This documents that Open guarantees connectivity, not just construction
	// of a sql.DB handle. A database which cannot be pinged must fail startup.
	t.Run("rejects a database that cannot be reached", func(t *testing.T) {
		cfg := testConfig(t)

		// A directory is not a valid sqlite database target.
		cfg.DB.ConnString = t.TempDir()

		_, err := db.Open(t.Context(), discardLogger(), &cfg.DB)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to open database")
	})
}

// TestClose verifies the ownership contract of DB.Close:
// once closed, the underlying SQL connection must no longer accept operations.
func TestClose(t *testing.T) {
	db, err := db.Open(t.Context(), discardLogger(), &testConfig(t).DB)
	require.NoError(t, err)

	require.NoError(t, db.Close())

	err = db.Connection.Ping()
	require.Error(t, err)

	// The exact driver error is not part of the package contract. The important
	// guarantee is that the closed connection cannot be reused.
	require.True(t, errors.Is(err, sql.ErrConnDone) || err != nil)
}
