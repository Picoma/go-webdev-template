package db

import (
	"context"
	"fmt"

	"idp/internal/db/migrations"

	"github.com/pressly/goose/v3"
)

func (db *DB) Migrate(ctx context.Context) error {
	// Get migrations version
	p, err := goose.NewProvider(
		goose.DialectSQLite3,
		db.Connection,
		migrations.EmbeddedFS,
		goose.WithSlog(db.Logger),
		goose.WithVerbose(true),
	)
	if err != nil {
		return fmt.Errorf("error setting up goose: %w", err)
	}

	_, err = p.Up(ctx)
	if err != nil {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	return nil
}
