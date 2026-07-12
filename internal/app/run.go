package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
)

// Run runs the app through its server.
// Logs wide events regarding the server status.
func (app *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	app.logger.InfoContext(ctx,
		"starting server",
		slog.String("server.address", app.server.Addr),
		slog.String("server.status", "open"),
	)

	go func() {
		errCh <- app.server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			msg := "server failed"
			app.logger.ErrorContext(ctx, msg,
				slog.String("server.address", app.server.Addr),
				slog.String("server.status", "crashed"),
				slog.String(app.cfg.LoggingSchema.ErrorType, "listen_error"),
				slog.Any(app.cfg.LoggingSchema.ErrorMessage, err),
			)
			return fmt.Errorf("%s: %w", msg, err)
		}

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			app.cfg.Server.ShutdownTimeout,
		)
		defer cancel()

		if err := app.server.Shutdown(shutdownCtx); err != nil {
			msg := "requested shutdown failed"
			app.logger.ErrorContext(ctx, msg,
				slog.String("server.address", app.server.Addr),
				slog.String("server.status", "crashed"),
				slog.String(app.cfg.LoggingSchema.ErrorType, "shutdown_error"),
				slog.Any(app.cfg.LoggingSchema.ErrorMessage, err),
			)
			return fmt.Errorf("%s: %w", msg, err)
		}
	}

	app.logger.InfoContext(ctx, "server stopped successfully",
		slog.String("server.address", app.server.Addr),
		slog.String("server.status", "stopped"),
	)
	return nil
}
