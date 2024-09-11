package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/repo"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/service"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/transport/http"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/pkg/httpserver"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/pkg/postgres"
)

func Run() int {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// app configuration
	cfg, err := NewConfig()
	if err != nil {
		logger.Error("failed to parse config", "err", err)
		return 1
	}

	// database connection
	pg, err := postgres.New(ctx, cfg.Postgres.Conn, postgres.DataTypes([]string{"organization_type", "tender_service_type", "tender_status", "bid_status", "bid_author_type", "bid_decision_type"}))
	if err != nil {
		logger.Error("failed to establish db conn", "err", err)
		return 1
	}
	defer func() {
		pg.Close()
		logger.Info("db conn closed", "uri", cfg.Postgres.Conn)
	}()
	logger.Info("db conn established", "uri", cfg.Postgres.Conn)

	// run migrations
	if code := Migrate(cfg, logger); code != 0 {
		return code
	}

	// repositories initialization
	employeeRepo := repo.NewEmployeePG(pg)
	tenderRepo := repo.NewTenderPG(pg)
	bidRepo := repo.NewBidPG(pg)
	decisionRepo := repo.NewBidDecisionPG(pg)
	reviewRepo := repo.NewBidReviewPG(pg)
	logger.Info("repositories initialized")

	// services initialization
	employeeService := service.NewEmployeeV1(employeeRepo)
	tenderService := service.NewTenderV1(tenderRepo, employeeService)
	bidService := service.NewBidV1(bidRepo, decisionRepo, tenderService, employeeService)
	reviewService := service.NewBidReviewV1(reviewRepo, bidService, tenderService, employeeService)
	logger.Info("services initialized")

	// http server start
	mux := http.NewMux(tenderService, bidService, reviewService, logger)
	server := httpserver.New(mux, httpserver.Addr(cfg.Server.Addr))
	server.Start(ctx)
	logger.Info("http server started", "addr", cfg.Server.Addr)

	// graceful shutdown
	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := server.Stop(ctx)
		if err != nil {
			logger.Error("failed to stop http server", "err", err)
		} else {
			logger.Info("http server has been stopped")
		}
		cancel()
		return 0
	case err := <-server.Err():
		logger.Error("http server returned error", "err", err)
		return 1
	}
}
