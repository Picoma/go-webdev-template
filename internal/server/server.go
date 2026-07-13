package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"idp/internal/config"
	"idp/internal/web"

	"github.com/go-chi/chi/v5"
)

func New(
	cfg *config.Server,
	router http.Handler,
) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.BindAddress, cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}

// NewRouter is decorrelated from [New] in order to facilitate tests,
// feeding the router directly to a [httptest.ResponseRecorder].
func NewRouter(
	logger *slog.Logger,
	cfg *config.Config,
	systemHandler SystemHandler,
	counterHandler CounterHandler,
) http.Handler {
	router := chi.NewRouter()

	registerMiddleware(router, logger)
	registerRoutes(router, cfg, systemHandler, counterHandler)
	web.RegisterStaticRoutes(router)

	return router
}
