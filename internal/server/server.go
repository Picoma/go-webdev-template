package server

import (
	"fmt"
	"net/http"

	"idp/internal/config"
)

func New(
	cfg config.Server,
	handler http.Handler,
) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.BindAddress, cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
