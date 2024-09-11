package app

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg *Config, logger *slog.Logger) {
	if cfg.Postgres.Migrations == "" {
		logger.Info("migrations path not found")
		return
	}

	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("%s?sslmode=disable", cfg.Postgres.Conn))

	if err != nil {
		logger.Error("failed to initialize migrations", "err", err)
		return
	}

	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("failed to apply migrations", "err", err)
		return
	}
}
