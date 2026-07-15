//nolint:paralleltest // If swapping databases is required/desired (like for Postgres), infra reuse might be needed, thus sequential tests
package db_test

// This file verifies the schema lifecycle contract of the database package.
//
// The application depends on migrations being:
//   - available from the embedded migration filesystem;
//   - applicable to a fresh database;
//   - safe to execute repeatedly.
//
// The tests intentionally validate the application's schema contract rather
// than goose internals. This keeps them useful even if the migration mechanism
// changes in the future.

import (
	"testing"

	"idp/internal/db"

	"github.com/stretchr/testify/require"
)

// TestMigrate verifies that a newly created database can be brought to the
// current schema version and that running migrations repeatedly is safe.
//
// Re-running migrations is important because deployments often execute startup
// migrations against databases that may already be partially or fully migrated.
func TestMigrate(t *testing.T) {
	ctx := t.Context()

	db, err := db.Open(ctx, discardLogger(), &testConfig(t).DB)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	t.Run("creates the application schema", func(t *testing.T) {
		require.NoError(t, db.Migrate(ctx))

		var tableName string
		err := db.Connection.QueryRowContext(
			ctx,
			`
			SELECT name
			FROM sqlite_master
			WHERE type = 'table'
			AND name = 'counter'
			`,
		).Scan(&tableName)

		require.NoError(t, err)
		require.Equal(t, "counter", tableName)
	})

	t.Run("can be executed more than once", func(t *testing.T) {
		require.NoError(t, db.Migrate(ctx))
		require.NoError(t, db.Migrate(ctx))
	})
}
