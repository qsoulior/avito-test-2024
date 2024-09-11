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
		logger.Warn("migrations path not found")
		return
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.Postgres.Migrations),
		fmt.Sprintf("%s?sslmode=disable", cfg.Postgres.Conn))

	if err != nil {
		logger.Error("failed to initialize migrations", "err", err)
		return
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("migrations have no changes")
			return
		}

		logger.Error("failed to apply migrations", "err", err)
		return
	}

	logger.Info("migrations applied")
}
