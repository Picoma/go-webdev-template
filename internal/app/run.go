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

	// TODO replace wide events
	app.logger = app.logger.With(
		slog.Any("server", map[string]any{
			"address": app.server.Addr,
			"status":  "open",
		}),
	)

	errCh := make(chan error, 1)

	app.logger.InfoContext(ctx, "starting server")

	go func() {
		errCh <- app.server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			msg := "server failed"

			app.logger.ErrorContext(ctx, msg,
				slog.Any("exception", map[string]any{
					"type":    "listen_error",
					"message": err,
				}),
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
				slog.Any("exception", map[string]any{
					"type":    "listen_error",
					"message": err,
				}),
			)

			return fmt.Errorf("%s: %w", msg, err)
		}
	}

	app.logger.InfoContext(ctx, "server stopped successfully")
	return nil
}
