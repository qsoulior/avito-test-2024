package http

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/service"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/transport/http/handler"
)

type Server struct {
	httpServer *http.Server
	errCh      chan error
}

func NewServer(addr string, tenderService service.Tender, bidService service.Bid, logger *slog.Logger) *Server {
	if logger == nil {
		return nil
	}

	router := http.NewServeMux()
	router.Handle("GET /ping", handler.Ping{})

	router.Handle("GET /tenders", handler.TenderGetByServiceType{Service: tenderService})
	router.Handle("POST /tenders/new", handler.TenderCreate{Service: tenderService})
	router.Handle("GET /tenders/my", handler.TenderGetByCreator{Service: tenderService})
	router.Handle("GET /tenders/{tenderId}/status", handler.TenderGetStatus{Service: tenderService})
	router.Handle("PUT /tenders/{tenderId}/status", handler.TenderUpdateStatus{Service: tenderService})
	router.Handle("PATCH /tenders/{tenderId}/edit", handler.TenderUpdate{Service: tenderService})
	router.Handle("PUT /tenders/{tenderId}/rollback/{version}", handler.TenderRollback{Service: tenderService})

	router.Handle("POST /bids/new", handler.BidCreate{Service: bidService})
	router.Handle("GET /bids/my", handler.BidGetByCreator{Service: bidService})
	router.Handle("GET /bids/{tenderId}/list", handler.BidGetByTender{Service: bidService})
	router.Handle("GET /bids/{bidId}/status", handler.BidGetStatus{Service: bidService})
	router.Handle("PUT /bids/{bidId}/status", handler.BidUpdateStatus{Service: bidService})
	router.Handle("PATCH /bids/{bidId}/edit", handler.BidUpdate{Service: bidService})
	router.Handle("PUT /bids/{bidId}/submit_decision", handler.BidSubmitDecision{Service: bidService})
	router.Handle("PUT /bids/{bidId}/rollback/{version}", handler.BidRollback{Service: bidService})

	var mux http.Handler = router
	middlewares := []Middleware{RecovererMiddleware(logger), LoggerMiddleware(logger)}
	for _, middleware := range middlewares {
		mux = middleware(mux)
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &Server{
		httpServer: httpServer,
		errCh:      make(chan error, 1),
	}
}

func (s *Server) Start(ctx context.Context) {
	s.httpServer.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	go func() {
		s.errCh <- s.httpServer.ListenAndServe()
		close(s.errCh)
	}()
}

func (s *Server) Stop(ctx context.Context) error { return s.httpServer.Shutdown(ctx) }

func (s *Server) Err() <-chan error { return s.errCh }
