package server

import (
	"net/http"

	"idp/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func registerRoutes(
	router chi.Router,
	cfg *config.Config,
	systemHandler SystemHandler,
	counterHandler CounterHandler,
) {
	router.Group(func(r chi.Router) {
		r.Get("/", http.RedirectHandler("/counter", http.StatusMovedPermanently).ServeHTTP)
		r.Get("/health", systemHandler.DisplayDBHealth)

		if cfg.WithDebug {
			r.Mount("/debug", middleware.Profiler())
		}
	})

	router.Group(func(r chi.Router) {
		// r.Use(AuthMiddleware)

		r.Get("/counter", counterHandler.DisplayCounter)
		r.Post("/counter", counterHandler.IncreaseCounter)
	})
}
