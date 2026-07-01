package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"idp/internal/config"
	"idp/internal/db"
	"idp/internal/handler"
	"idp/internal/server"
	"idp/internal/service"
)

type App struct {
	logger *slog.Logger
	cfg    *config.Config

	db     *db.DB
	server *http.Server
}

// New initializes the app (creating db connection, configuring server).
func New(ctx context.Context, logger *slog.Logger, cfg *config.Config) (*App, error) {
	database, err := db.Open(ctx, logger, cfg)
	if err != nil {
		msg := "error opening database"
		logger.ErrorContext(ctx, msg,
			slog.Any("exception.message", err),
		)
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	counterSvc := service.NewCountService(database.Queries)
	counterHandler := handler.NewCounterHandler(counterSvc)

	router := server.NewRouter(
		logger,
		cfg,
		database,
		counterHandler,
	)

	srv := server.New(
		cfg.Server,
		router,
	)

	return &App{
		logger: logger,
		cfg:    cfg,

		db:     database,
		server: srv,
	}, nil
}

func (app *App) Close() error {
	if err := app.db.Close(); err != nil {
		return fmt.Errorf("error closing application: %w", err)
	}

	return nil
}
