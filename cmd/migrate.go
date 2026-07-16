package cmd

import (
	"context"
	"fmt"

	"idp/internal/config"
	"idp/internal/db"
	"idp/internal/log"

	"github.com/urfave/cli/v3"
)

func newMigrateCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "Run database migrations",
		Action: func(ctx context.Context, c *cli.Command) error {
			logger := log.New(c.Writer, cfg)

			db, err := db.Open(ctx, logger, &cfg.DB)
			if err != nil {
				return fmt.Errorf("error connection to database: %w", err)
			}

			err = db.Migrate(ctx)
			if err != nil {
				return fmt.Errorf("error running migrations: %w", err)
			}

			return nil
		},
	}
}
