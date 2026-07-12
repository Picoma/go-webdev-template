package config

import (
	"time"

	"github.com/go-chi/httplog/v3"
)

type Service struct {
	Name        string
	Description string
	Version     string
	Env         string
	Commit      string
}

type DB struct {
	ConnString                     string
	Driver                         string
	HealthyOpenConnectionThreshold int
	HealthyWaitCountThreshold      int
}

type Server struct {
	BindAddress     string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type Config struct {
	Service Service
	DB      DB
	Server  Server

	Debug         bool
	TintedLogs    bool
	LoggingSchema *httplog.Schema
}

func Defaults(service Service) *Config {
	//nolint:mnd // default config
	return &Config{
		Service: service,
		DB: DB{
			ConnString:                     "", // Set in CLI
			Driver:                         "sqlite3",
			HealthyOpenConnectionThreshold: 40,
			HealthyWaitCountThreshold:      1000,
		},
		Server: Server{
			BindAddress:     "", // Set in CLI
			Port:            0,  // Set in CLI
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     time.Minute,
			ShutdownTimeout: 30 * time.Second,
		},
		Debug:         false, // Set in CLI
		TintedLogs:    false,
		LoggingSchema: httplog.SchemaOTEL,
	}
}
