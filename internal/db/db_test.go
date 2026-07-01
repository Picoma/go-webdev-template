package db_test

import (
	"context"
	"log/slog"
	"testing"

	"idp/internal/config"
	"idp/internal/db"
)

func TestOpen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  string
		cfg  *config.Config
	}{
		{
			name: "sqlite memory database",
			url:  ":memory:",
			cfg: &config.Config{
				DB: config.DB{
					Driver:     "sqlite3",
					ConnString: ":memory:",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := slog.Default()

			database, err := db.Open(context.Background(), logger, tt.cfg)
			if err != nil {
				t.Fatalf("Open() error = %v", err)
			}

			if database == nil {
				t.Fatal("expected database instance")
			}

			if database.Connection == nil {
				t.Fatal("expected sql connection")
			}

			if database.Queries == nil {
				t.Fatal("expected queries")
			}

			if err := database.Close(); err != nil {
				t.Fatalf("Close() error = %v", err)
			}
		})
	}
}

func TestDB_Close(t *testing.T) {
	t.Parallel()

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

	if err := database.Connection.Ping(); err == nil {
		t.Fatal("expected ping to fail after close")
	}
}
