package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"idp/internal/config"
	"idp/internal/db/migrations"
	"idp/internal/db/queries"

	_ "github.com/mattn/go-sqlite3" // SQLite drivers
	"github.com/pressly/goose/v3"
)

//go:generate sqlc generate -f ../../sqlc.yaml

// DB being the entrypoint for services, he implements [server.HealthChecker].
type DB struct {
	Logger *slog.Logger
	cfg    *config.DB

	Connection *sql.DB
	Queries    *queries.Queries
}

func Open(ctx context.Context, logger *slog.Logger, cfg *config.DB) (*DB, error) {
	logger = logger.WithGroup("db").With(
		slog.String("server.address", cfg.ConnString),
		slog.String("system.name", cfg.Driver),
	)

	// Open connection
	sqlDB, err := sql.Open(cfg.Driver, cfg.ConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Check connection works
	// if doesn't, gracefully close connection
	success := false
	defer func() {
		if !success {
			_ = sqlDB.Close()
		}
	}()

	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Get migrations version
	p, err := goose.NewProvider(
		goose.DialectSQLite3,
		sqlDB,
		migrations.EmbeddedFS,
		goose.WithSlog(logger),
		goose.WithVerbose(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error setting up goose: %w", err)
	}

	// Checks pending migrations, warns if has pending
	currentVersion, targetVersion, err := p.GetVersions(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking pending migrations: %w", err)
	}
	if currentVersion != targetVersion {
		msg := "pending migrations ; app may behave incorrectly."
		logger.WarnContext(ctx, msg,
			slog.Int64("db.schema.current", currentVersion),
			slog.Int64("db.schema.target", targetVersion),
		)
	}

	success = true
	logger.InfoContext(ctx, "database connection successful")
	db := &DB{
		Logger: logger,
		cfg:    cfg,

		Connection: sqlDB,
		Queries:    queries.New(sqlDB),
	}

	return db, nil
}

func (db *DB) Close() error {
	if err := db.Connection.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	db.Logger.Info("database connection closed")
	return nil
}
