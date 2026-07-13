package server_test

import (
	"net/http"
)

type fakeCounterHandler struct {
	displayCalled   bool
	incrementCalled bool
}

func (h *fakeCounterHandler) DisplayCounter(w http.ResponseWriter, _ *http.Request) {
	h.displayCalled = true
	w.WriteHeader(http.StatusNoContent)
}

func (h *fakeCounterHandler) IncreaseCounter(w http.ResponseWriter, _ *http.Request) {
	h.incrementCalled = true
	w.WriteHeader(http.StatusNoContent)
}

type fakeSystemHandler struct {
	healthCalled bool
}

func (h *fakeSystemHandler) DisplayDBHealth(w http.ResponseWriter, _ *http.Request) {
	h.healthCalled = true
	w.WriteHeader(http.StatusNoContent)
}
