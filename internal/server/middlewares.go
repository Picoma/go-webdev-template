package server

import (
	"log/slog"
	"net/http"

	"idp/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
)

//nolint:mnd // Package-wide configuration options, not exposed to final user
func registerMiddleware(r chi.Router, logger *slog.Logger, cfg *config.Config) {
	loggingOptions := &httplog.Options{
		Level:         slog.LevelInfo,
		Schema:        cfg.LoggingSchema,
		RecoverPanics: true,
		Skip: func(_ *http.Request, respStatus int) bool {
			return respStatus == 404 || respStatus == 405
		},

		LogRequestHeaders:  []string{"Origin"},
		LogResponseHeaders: []string{},
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
	r.Use(httplog.RequestLogger(logger, loggingOptions))
}
