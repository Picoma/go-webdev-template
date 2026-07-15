package schema

import (
	"log/slog"

	"idp/internal/config"
)

func DB(driver string, connString string) slog.Attr {
	return slog.Group("db",
		slog.Group("system",
			slog.String("name", driver),
		),
		slog.Group("server",
			slog.String("address", connString),
		),
	)
}

func DBSchema(current *int64, target *int64) slog.Attr {
	return slog.Group("db",
		slog.Group("schema",
			slog.Int64("current", *current),
			slog.Int64("target", *target),
		),
	)
}

func Service(service config.Service) slog.Attr {
	return slog.Group("service",
		slog.String("name", service.Name),
		slog.String("version", service.Version),
		slog.String("hash_commit", service.HashCommit),
		slog.String("env", service.Env),
	)
}
