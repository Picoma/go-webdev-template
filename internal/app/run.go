package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"idp/internal/log"
)

// Run runs the app through its server.
// Logs wide events regarding the server status.
func (app *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.AddContextAttrs(ctx,
		slog.String("server.address", app.server.Addr),
		slog.String("server.status", "open"),
	)

	errCh := make(chan error, 1)

	app.logger.InfoContext(ctx, "starting server")

	go func() {
		errCh <- app.server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed: %w", err)
		}

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			app.cfg.Server.ShutdownTimeout,
		)
		defer cancel()

		if err := app.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("requested shutdown failed: %w", err)
		}
	}

	app.logger.InfoContext(ctx, "server stopped successfully")
	return nil
}
