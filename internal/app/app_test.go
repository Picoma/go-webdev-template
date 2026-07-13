package app_test

import (
	"context"
	"log/slog"
	"testing"

	"idp/internal/app"
	"idp/internal/config"
)

func TestNew(t *testing.T) {
	t.Parallel()

	app, err := app.New(context.Background(), slog.Default(), &config.Config{
		Service: config.Service{
			Version: "test",
		},
		DB: config.DB{
			Driver:     "sqlite3",
			ConnString: ":memory:",
		},
		Server:     config.Server{Port: 8080},
		TintedLogs: false,
		Verbose:    false,
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Cleanup(func() {
		if err := app.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	})
}

func TestAppClose(t *testing.T) {
	t.Parallel()

	app, err := app.New(context.Background(), slog.Default(), &config.Config{
		DB: config.DB{
			Driver:     "sqlite3",
			ConnString: ":memory:",
		},
		Server:     config.Server{Port: 8080},
		Service:    config.Service{Version: "test"},
		TintedLogs: false,
		Verbose:    false,
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := app.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
