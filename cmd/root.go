package cmd

import (
	"idp/internal/config"

	"github.com/urfave/cli/v3"
)

func NewRootCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:      cfg.Service.Name,
		Version:   cfg.Service.Version,
		Usage:     cfg.Service.Description,
		Authors:   []any{},
		Copyright: "",
		Metadata:  map[string]any{},

		Commands: []*cli.Command{
			newServeCmd(cfg),
			newMigrateCmd(cfg),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "format-logs",
				Usage:       "log events in human readable format",
				Aliases:     []string{"f"},
				Sources:     cli.EnvVars("IDP_FORMAT_LOGS"),
				Destination: &cfg.TintedLogs,
				Value:       false,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "verbose output",
				Sources:     cli.EnvVars("IDP_VERBOSE"),
				Destination: &cfg.Verbose,
				Value:       false,
			},
			&cli.StringFlag{
				Name:        "db-string",
				Usage:       "database connection string",
				Sources:     cli.EnvVars("IDP_DB_STRING"),
				Destination: &cfg.DB.ConnString,
				Value:       "idp.sqlite",
			},
		},

		EnableShellCompletion: false,

		Reader:    nil,
		Writer:    nil,
		ErrWriter: nil,

		Before: nil,
		After:  nil,
		Action: nil,
	}
}
