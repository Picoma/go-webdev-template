package cmd

import (
	"context"
	"fmt"

	"idp/internal/app"
	"idp/internal/config"
	"idp/internal/log"

	"github.com/urfave/cli/v3"
)

func newServeCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Start the application server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Usage:       "server binding address",
				Aliases:     []string{"a"},
				Sources:     cli.EnvVars("IDP_ADDRESS"),
				Destination: &cfg.Server.BindAddress,
				Value:       "0.0.0.0",
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "server port to listen on",
				Aliases:     []string{"p"},
				Sources:     cli.EnvVars("IDP_PORT"),
				Destination: &cfg.Server.Port,
				Value:       8080,
			},
			&cli.BoolFlag{
				Name:        "with-debug",
				Usage:       "exposes a profiling endpoint on /debug",
				Aliases:     []string{"d"},
				Sources:     cli.EnvVars("IDP_DEBUG"),
				Destination: &cfg.WithDebug,
				Value:       false,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			ctx, logger := log.NewWithContext(ctx, c.Writer, cfg)

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
