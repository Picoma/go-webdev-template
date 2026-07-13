package server_test

// This file documents the contract of the HTTP server constructor.
//
// The New function is intentionally a thin configuration adapter. It does not
// start the server and does not perform networking; it only translates service
// configuration into an http.Server instance.
//
// Tests therefore focus exclusively on configuration propagation.

import (
	"net/http"
	"testing"
	"time"

	"idp/internal/config"
	"idp/internal/server"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	router := http.NewServeMux()

	cfg := &config.Server{
		BindAddress:  "127.0.0.1",
		Port:         8080,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	srv := server.New(cfg, router)

	require.Equal(t, "127.0.0.1:8080", srv.Addr)
	require.Same(t, router, srv.Handler)

	require.Equal(t, cfg.ReadTimeout, srv.ReadTimeout)
	require.Equal(t, cfg.WriteTimeout, srv.WriteTimeout)
	require.Equal(t, cfg.IdleTimeout, srv.IdleTimeout)
}
