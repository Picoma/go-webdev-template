package cmd

import (
	"context"
	"fmt"

	"idp/internal/app"
	"idp/internal/config"
	"idp/internal/log"

	"github.com/urfave/cli/v3"
)

func newRunCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Start the application server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "bind-address",
				Usage:       "Server address to bind to",
				Aliases:     []string{"b"},
				Destination: &cfg.Server.BindAddress,
				Value:       "0.0.0.0",
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Server port to listen on",
				Aliases:     []string{"p"},
				Destination: &cfg.Server.Port,
				Value:       8080,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			logger := log.New(c.Writer, cfg)

			application, err := app.New(ctx, logger, cfg)
			if err != nil {
				return fmt.Errorf("error building application: %w", err)
			}
			defer application.Close()

			if err := application.Run(ctx); err != nil {
				return fmt.Errorf("failed to start application: %w", err)
			}

			return nil
		},
	}
}
