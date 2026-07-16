package server

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	slogchi "github.com/samber/slog-chi"
)

//nolint:mnd // Package-wide configuration options, not exposed to final user
func registerMiddleware(r chi.Router, logger *slog.Logger) {
	loggingOptions := slogchi.Config{
		WithSpanID:  true,
		WithTraceID: true,
	}

	corsOptions := cors.Options{
		AllowedOrigins: []string{
			"https://*",
			"http://*",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
		},
		AllowCredentials: false,
		MaxAge:           300,
	}

	r.Use(cors.Handler(corsOptions))
	r.Use(slogchi.NewWithConfig(logger.WithGroup("http"), loggingOptions))
}
