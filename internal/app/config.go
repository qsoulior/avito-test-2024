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

// Config
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

// ConfigServer
type ConfigServer struct {
	Addr string
}

const (
	SERVER_ADDRESS = "SERVER_ADDRESS"
)

func (c *ConfigServer) ParseEnv() error {
	addr, ok := os.LookupEnv(SERVER_ADDRESS)
	if !ok {
		return NewEnvError(SERVER_ADDRESS)
	}

	c.Addr = addr
	return nil
}

// ConfigPostgres
type ConfigPostgres struct {
	Conn string
}

const (
	POSTGRES_CONN     = "POSTGRES_CONN"
	POSTGRES_USERNAME = "POSTGRES_USERNAME"
	POSTGRES_PASSWORD = "POSTGRES_PASSWORD"
	POSTGRES_HOST     = "POSTGRES_HOST"
	POSTGRES_PORT     = "POSTGRES_PORT"
	POSTGRES_DATABASE = "POSTGRES_DATABASE"
)

func (c *ConfigPostgres) ParseEnv() error {
	conn, ok := os.LookupEnv(POSTGRES_CONN)
	if ok {
		c.Conn = conn
		return nil
	}

	username, ok := os.LookupEnv(POSTGRES_USERNAME)
	if !ok {
		return NewEnvError(POSTGRES_USERNAME)
	}

	password, ok := os.LookupEnv(POSTGRES_PASSWORD)
	if !ok {
		return NewEnvError(POSTGRES_PASSWORD)
	}

	host, ok := os.LookupEnv(POSTGRES_HOST)
	if !ok {
		return NewEnvError(POSTGRES_HOST)
	}

	port, ok := os.LookupEnv(POSTGRES_PORT)
	if !ok {
		return NewEnvError(POSTGRES_PORT)
	}

	database, ok := os.LookupEnv(POSTGRES_DATABASE)
	if !ok {
		return NewEnvError(POSTGRES_DATABASE)
	}

	c.Conn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database)
	return nil
}

// NewConfig
func NewConfig() (*Config, error) {
	cfg := new(Config)
	err := cfg.ParseEnv()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
