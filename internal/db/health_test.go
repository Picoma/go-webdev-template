//nolint:paralleltest // If swapping databases is required/desired (like for Postgres), infra reuse might be needed, thus sequential tests
package db_test

// This file documents the operational contract exposed by CheckDatabase.
//
// A health check is expected to:
//   - report a reachable database as healthy;
//   - report a broken connection as unavailable.
//
// The tests avoid forcing connection-pool statistics because those values are
// runtime characteristics rather than stable behavioral guarantees. Testing
// them would make the suite dependent on database/sql scheduling behavior.

import (
	"testing"

	"idp/internal/db"

	"github.com/stretchr/testify/require"
)

// TestCheckDatabase verifies the two externally meaningful states of the
// database health endpoint: healthy and unavailable.
func TestCheckDatabase(t *testing.T) {
	ctx := t.Context()

	t.Run("reports a healthy database", func(t *testing.T) {
		db, err := db.Open(ctx, discardLogger(), testConfig(t))
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, db.Close())
		})

		status := db.CheckDatabase(ctx)

		require.Equal(t, "up", status.Status)
		require.Equal(t, "healthy", status.Message)
		require.Empty(t, status.Error)
	})

	// This protects the failure path used by monitoring systems. A database
	// which cannot answer a ping must never be reported as healthy.
	t.Run("reports an unavailable database", func(t *testing.T) {
		db, err := db.Open(ctx, discardLogger(), testConfig(t))
		require.NoError(t, err)

		require.NoError(t, db.Connection.Close())

		status := db.CheckDatabase(ctx)

		require.Equal(t, "down", status.Status)
		require.NotEmpty(t, status.Error)
	})
}
