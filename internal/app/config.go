package app

import (
	"fmt"
	"os"
)

type EnvError struct {
	env string
}

func NewEnvError(env string) error { return &EnvError{env} }

func (e *EnvError) Error() string { return fmt.Sprintf("could not find env variable: %s", e.env) }

// Config.
type Config struct {
	Server   ConfigServer
	Postgres ConfigPostgres
}

func (c *Config) ParseEnv() error {
	if err := c.Server.ParseEnv(); err != nil {
		return err
	}

	return c.Postgres.ParseEnv()
}

// ConfigServer.
type ConfigServer struct {
	Addr string
}

const (
	EnvServerAddress = "SERVER_ADDRESS"
)

func (c *ConfigServer) ParseEnv() error {
	addr, ok := os.LookupEnv(EnvServerAddress)
	if !ok {
		return NewEnvError(EnvServerAddress)
	}

	c.Addr = addr
	return nil
}

// ConfigPostgres.
type ConfigPostgres struct {
	Conn       string
	Migrations string
}

const (
	EnvPostgresConn       = "POSTGRES_CONN"
	EnvPostgresUsername   = "POSTGRES_USERNAME"
	EnvPostgresPassword   = "POSTGRES_PASSWORD"
	EnvPostgresHost       = "POSTGRES_HOST"
	EnvPostgresPort       = "POSTGRES_PORT"
	EnvPostgresDatabase   = "POSTGRES_DATABASE"
	EnvPostgresMigrations = "POSTGRES_MIGRATIONS"
)

func (c *ConfigPostgres) ParseEnv() error {
	c.Migrations = os.Getenv(EnvPostgresMigrations)

	conn, ok := os.LookupEnv(EnvPostgresConn)
	if ok {
		c.Conn = conn
		return nil
	}

	username, ok := os.LookupEnv(EnvPostgresUsername)
	if !ok {
		return NewEnvError(EnvPostgresUsername)
	}

	password, ok := os.LookupEnv(EnvPostgresPassword)
	if !ok {
		return NewEnvError(EnvPostgresPassword)
	}

	host, ok := os.LookupEnv(EnvPostgresHost)
	if !ok {
		return NewEnvError(EnvPostgresHost)
	}

	port, ok := os.LookupEnv(EnvPostgresPort)
	if !ok {
		return NewEnvError(EnvPostgresPort)
	}

	database, ok := os.LookupEnv(EnvPostgresDatabase)
	if !ok {
		return NewEnvError(EnvPostgresDatabase)
	}

	c.Conn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database)

	return nil
}

// NewConfig.
func NewConfig() (*Config, error) {
	cfg := new(Config)
	err := cfg.ParseEnv()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
