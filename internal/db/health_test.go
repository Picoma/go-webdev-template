package db_test

import (
	"context"
	"log/slog"
	"testing"

	"idp/internal/config"
	"idp/internal/db"
)

func TestDB_CheckDatabase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T) *db.DB
		wantStatus string
		wantMsg    string
	}{
		{
			name: "healthy database",
			setup: func(t *testing.T) *db.DB {
				t.Helper()

				database, err := db.Open(
					context.Background(),
					slog.Default(),
					&config.Config{
						DB: config.DB{
							Driver:                         "sqlite3",
							ConnString:                     ":memory:",
							HealthyOpenConnectionThreshold: 40,
							HealthyWaitCountThreshold:      1000,
						},
					},
				)
				if err != nil {
					t.Fatalf("Open() error = %v", err)
				}

				t.Cleanup(func() {
					_ = database.Close()
				})

				return database
			},
			wantStatus: "up",
			wantMsg:    "healthy",
		},
		{
			name: "closed database",
			setup: func(t *testing.T) *db.DB {
				t.Helper()

				database, err := db.Open(
					context.Background(),
					slog.Default(),
					&config.Config{
						DB: config.DB{
							Driver:     "sqlite3",
							ConnString: ":memory:",
						},
					},
				)
				if err != nil {
					t.Fatalf("Open() error = %v", err)
				}

				if err := database.Close(); err != nil {
					t.Fatalf("Close() error = %v", err)
				}

				return database
			},
			wantStatus: "down",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			database := tt.setup(t)

			status := database.CheckDatabase(
				context.Background(),
			)

			if status.Status != tt.wantStatus {
				t.Fatalf(
					"Status: got %q want %q",
					status.Status,
					tt.wantStatus,
				)
			}

			if status.Message != tt.wantMsg {
				t.Fatalf(
					"Message: got %q want %q",
					status.Message,
					tt.wantMsg,
				)
			}
		})
	}
}
