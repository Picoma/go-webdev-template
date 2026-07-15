package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func registerRoutes(
	router chi.Router,
	debug bool,
	systemHandler SystemHandler,
	counterHandler CounterHandler,
) {
	router.Group(func(r chi.Router) {
		r.Get("/", http.RedirectHandler("/counter", http.StatusMovedPermanently).ServeHTTP)
		r.Get("/health", systemHandler.DisplayDBHealth)

		if debug {
			r.Mount("/debug", middleware.Profiler())
		}
	})

	router.Group(func(r chi.Router) {
		// r.Use(AuthMiddleware)

		r.Get("/counter", counterHandler.DisplayCounter)
		r.Post("/counter", counterHandler.IncreaseCounter)
	})
}
