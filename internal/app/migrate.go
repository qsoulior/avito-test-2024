package app

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg *Config, logger *slog.Logger) int {
	if cfg.Postgres.Migrations == "" {
		logger.Warn("migrations path not found")
		return 0
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.Postgres.Migrations),
		cfg.Postgres.Conn)

	if err != nil {
		logger.Error("failed to initialize migrations", "err", err)
		return 1
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("migrations have no changes")
			return 0
		}

		logger.Error("failed to apply migrations", "err", err)
		return 1
	}

	logger.Info("migrations applied")
	return 0
}
