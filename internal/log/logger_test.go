package log_test

// This file documents the logging contract of the application.
//
// The log package is infrastructure code: its output is consumed by humans
// during development and by log aggregation systems in production. Tests
// therefore focus on stable observable behavior:
//
//   - configured log levels are respected;
//   - structured service metadata is always attached;
//   - machine-readable output remains valid JSON;
//   - the global slog default is updated.
//
// The tests intentionally avoid asserting formatting details such as colors,
// timestamps, or exact text layout. Those belong to the underlying slog handler
// implementation and may evolve independently.

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"idp/internal/config"
	"idp/internal/log"

	"github.com/stretchr/testify/require"
)

func testConfig() *config.Config {
	service := config.Service{
		Name:       "idp_test",
		Version:    "1.2.3",
		HashCommit: "abcdef",
		Env:        "test",
	}

	cfg := config.Defaults(service)

	return cfg
}

func TestNew_JSONLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		configure func(*config.Config)
		log       func(context.Context, *slog.Logger)
		assert    func(t *testing.T, record map[string]any)
	}{
		{
			name: "attaches service metadata",
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "startup")
			},
			assert: func(t *testing.T, record map[string]any) {
				require.Contains(t,
					record,
					"service.name",
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var output bytes.Buffer

			cfg := testConfig()
			if tt.configure != nil {
				tt.configure(cfg)
			}

			ctx, logger := log.NewWithContext(t.Context(), &output, cfg)
			tt.log(ctx, logger)

			var record map[string]any

			require.NoError(t,
				json.Unmarshal(output.Bytes(), &record),
			)

			tt.assert(t, record)
		})
	}
}

func TestNewLogLevelFiltering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		verbose    bool
		log        func(context.Context, *slog.Logger)
		wantOutput bool
	}{
		{
			name:    "debug hidden when verbose disabled",
			verbose: false,
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.DebugContext(ctx, "debug message")
			},
			wantOutput: false,
		},
		{
			name:    "debug visible when verbose enabled",
			verbose: true,
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.DebugContext(ctx, "debug message")
			},
			wantOutput: true,
		},
		{
			name:    "info always visible",
			verbose: false,
			log: func(ctx context.Context, logger *slog.Logger) {
				logger.InfoContext(ctx, "info message")
			},
			wantOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := testConfig()
			cfg.Verbose = tt.verbose

			var output bytes.Buffer

			ctx, logger := log.NewWithContext(t.Context(), &output, cfg)

			tt.log(ctx, logger)

			if tt.wantOutput {
				require.NotEmpty(t, output.String())
			} else {
				require.Empty(t, output.String())
			}
		})
	}
}
