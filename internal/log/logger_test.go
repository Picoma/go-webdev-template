package log_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"idp/internal/config"
	"idp/internal/log"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config *config.Config
		want   slog.Level
	}{
		{
			name: "default info level",
			config: &config.Config{
				Service: config.Service{
					Version: "test",
					Env:     "test",
				},
			},
			want: slog.LevelInfo,
		},
		{
			name: "verbose debug level",
			config: &config.Config{
				Service: config.Service{
					Version: "test",
					Env:     "test",
				},
				Debug: true,
			},
			want: slog.LevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := log.New(os.Stdout, tt.config)

			if logger == nil {
				t.Fatal("expected logger")
			}

			if !logger.Enabled(
				context.Background(),
				tt.want,
			) {
				t.Fatalf(
					"expected level %v enabled",
					tt.want,
				)
			}
		})
	}
}
