package db

import (
	"context"
	"time"
)

type DatabaseStatus struct {
	Status            string        `json:"status"`
	Message           string        `json:"message"`
	Error             string        `json:"error"`
	OpenConnections   int           `json:"open_connections"`
	InUse             int           `json:"in_use"`
	Idle              int           `json:"idle"`
	WaitCount         int64         `json:"wait_count"`
	WaitDuration      time.Duration `json:"wait_duration"`
	MaxIdleClosed     int64         `json:"max_idle_closed"`
	MaxLifetimeClosed int64         `json:"max_lifetime_closed"`
}

func (db *DB) CheckDatabase(ctx context.Context) DatabaseStatus {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := db.Connection.PingContext(ctx); err != nil {
		return DatabaseStatus{
			Status: "down",
			Error:  err.Error(),
		}
	}

	stats := db.Connection.Stats()

	out := DatabaseStatus{
		Status:            "up",
		Message:           "healthy",
		OpenConnections:   stats.OpenConnections,
		InUse:             stats.InUse,
		Idle:              stats.Idle,
		WaitCount:         stats.WaitCount,
		WaitDuration:      stats.WaitDuration,
		MaxIdleClosed:     stats.MaxIdleClosed,
		MaxLifetimeClosed: stats.MaxLifetimeClosed,
	}

	// thresholds
	if stats.OpenConnections > db.cfg.DB.HealthyOpenConnectionThreshold {
		out.Message = "high load"
	}

	if stats.WaitCount > int64(db.cfg.DB.HealthyWaitCountThreshold) {
		out.Message = "contention detected"
	}

	return out
}
