package config

import (
	"time"

	"github.com/go-chi/httplog/v3"
)

// Config holds the config schema used throuout the code.
//
// Instanciate using [Defaults].
// JSON tags present to serialize in structured logging.
type Config struct {
	Service Service `json:"-"`
	DB      DB      `json:"db"`
	Server  Server  `json:"server"`

	WithDebug     bool            `json:"with_debug"`
	Verbose       bool            `json:"verbose"`
	LoggingSchema *httplog.Schema `json:"-"`
}

type Service struct {
	Name        string `json:"name"`
	Description string `json:"-"`
	Version     string `json:"version"`
	Env         string `json:"env"`
	HashCommit  string `json:"hash_commit"`
}

type DB struct {
	ConnString                     string `json:"conn_string"`
	Driver                         string `json:"driver"`
	HealthyOpenConnectionThreshold int    `json:"healthy_open_connection_threshold"`
	HealthyWaitCountThreshold      int    `json:"healthy_wait_count_threshold"`
}

type Server struct {
	BindAddress     string        `json:"bind_address"`
	Port            int           `json:"port"`
	ReadTimeout     time.Duration `json:"read_timeout_ns"`
	WriteTimeout    time.Duration `json:"write_timeout_ns"`
	IdleTimeout     time.Duration `json:"idle_timeout_ns"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout_ns"`
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
		WithDebug:     false, // Set in CLI
		Verbose:       false, // Set in CLI
		LoggingSchema: httplog.SchemaOTEL,
	}
}
