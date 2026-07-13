package handler

import (
	"encoding/json"
	"net/http"

	"idp/internal/db"
)

// SystemHandler implements [server.SystemHandler].
//
// It uses the database directly, thus merging the service and handler layers.
// This is fine as ther will never be "business" logic here : we display the raw
// system data.
type SystemHandler struct {
	db *db.DB
}

func NewSystemHandler(db *db.DB) *SystemHandler {
	return &SystemHandler{
		db: db,
	}
}

func (sh *SystemHandler) DisplayDBHealth(w http.ResponseWriter, r *http.Request) {
	status := sh.db.CheckDatabase(r.Context())

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

// // PromMetrics displays a prometheus metrics endpoint
// func (sh *SystemHandler) PromMetrics(w http.ResponseWriter, r *http.Request) {}

// // DisplayVersion displays the current version, queriable from HTTP.
// func (sh *SystemHandler) DisplayVersion(w http.ResponseWriter, r *http.Request) {}
