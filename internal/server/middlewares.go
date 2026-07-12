package server

import (
	"log/slog"
	"net/http"

	"idp/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
)

//nolint:mnd // Package-wide configuration options, not exposed to final user
func registerMiddleware(r chi.Router, cfg *config.Config, logger *slog.Logger) {
	loggingOptions := httplog.Options{
		Level:         slog.LevelInfo,
		Schema:        cfg.LoggingSchema,
		RecoverPanics: true,
		Skip: func(_ *http.Request, status int) bool {
			return status == 404 || status == 405
		},
		LogRequestHeaders: []string{
			"Accept-Encoding",
			"Content-Type",
			"Content-Length",
			"Connection",
			"Host",
			"HX-Request",
			"Origin",
			"User-agent",
			"X-Forwarded-For",
			"X-Forwarded-Proto",
			"X-Real-Ip",
		},
		LogExtraAttrs: func(req *http.Request, _ string, _ int) []slog.Attr {
			reqID := middleware.GetReqID(req.Context())
			return []slog.Attr{
				slog.String("request.id", reqID),
			}
		},
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
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}

	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger, &loggingOptions))
	r.Use(cors.Handler(corsOptions))
}
