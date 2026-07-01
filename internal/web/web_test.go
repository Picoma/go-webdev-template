package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"idp/internal/web"

	"github.com/go-chi/chi/v5"
)

func TestRegisterStaticRoutes(t *testing.T) {
	t.Parallel()

	router := chi.NewRouter()

	web.RegisterStaticRoutes(router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/assets/",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Fatal("expected assets route to be registered")
	}
}
