package handler

import (
	"context"
	"encoding/json"

	"idp/internal/db"

	"net/http"
)

type HealthChecker interface {
	CheckDatabase(context.Context) db.DatabaseStatus
}

func HealthHandler(database HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := database.CheckDatabase(r.Context())

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
	}
}
