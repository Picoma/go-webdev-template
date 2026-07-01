package server

import (
	"log/slog"
	"net/http"

	"idp/internal/config"
	"idp/internal/handler"
	"idp/internal/web"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	logger *slog.Logger,
	cfg *config.Config,
	db handler.HealthChecker,
	ch *handler.CounterHandler,
) http.Handler {
	r := chi.NewRouter()

	registerMiddleware(r, cfg, logger)
	registerRoutes(r, cfg, db, ch)
	web.RegisterStaticRoutes(r)

	return r
}
