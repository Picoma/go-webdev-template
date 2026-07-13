package server

import (
	"net/http"

	"idp/internal/config"
	"idp/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func registerRoutes(
	r chi.Router,
	cfg *config.Config,
	db handler.HealthChecker,
	ch *handler.CounterHandler,
) {
	r.Group(func(r chi.Router) {
		r.Get("/", http.RedirectHandler("/counter", http.StatusMovedPermanently).ServeHTTP)
		r.Get("/health", handler.HealthHandler(db))
		if cfg.WithDebug {
			r.Mount("/debug", middleware.Profiler())
		}
	})

	r.Group(func(r chi.Router) {
		// r.Use(AuthMiddleware)

		r.Get("/counter", ch.DisplayCounter)
		r.Post("/counter", ch.IncreaseCounter)
	})
}
