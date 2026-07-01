package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"idp/internal/db"
	"idp/internal/handler"
)

type fakeHealthChecker struct {
	status db.DatabaseStatus
}

func (f fakeHealthChecker) CheckDatabase(context.Context) db.DatabaseStatus {
	return f.status
}

func TestHealthHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		status     db.DatabaseStatus
		wantStatus int
	}{
		{
			name: "healthy database",
			status: db.DatabaseStatus{
				Status:  "up",
				Message: "healthy",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "unhealthy database",
			status: db.DatabaseStatus{
				Status: "down",
				Error:  "connection refused",
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := handler.HealthHandler(fakeHealthChecker{
				status: tt.status,
			})

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d want %d", rec.Code, tt.wantStatus)
			}

			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf(
					"content type: got %q want %q",
					got,
					"application/json",
				)
			}

			var got db.DatabaseStatus
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if got.Status != tt.status.Status {
				t.Fatalf(
					"status field: got %q want %q",
					got.Status,
					tt.status.Status,
				)
			}

			if got.Message != tt.status.Message {
				t.Fatalf(
					"message field: got %q want %q",
					got.Message,
					tt.status.Message,
				)
			}
		})
	}
}
