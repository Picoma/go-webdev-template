package server

import (
	"net/http"
)

type CounterHandler interface {
	DisplayCounter(w http.ResponseWriter, r *http.Request)
	IncreaseCounter(w http.ResponseWriter, r *http.Request)
}

type SystemHandler interface {
	DisplayDBHealth(w http.ResponseWriter, r *http.Request)
}
