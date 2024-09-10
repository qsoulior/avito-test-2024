package http

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/transport/http/handler"
)

type Server struct {
	httpServer *http.Server
	errCh      chan error
}

func NewServer(host string, port string, logger *slog.Logger) *Server {
	if logger == nil {
		return nil
	}

	router := http.NewServeMux()
	router.Handle("GET /ping", handler.Ping{})

	var mux http.Handler = router
	middlewares := []Middleware{RecovererMiddleware(logger), LoggerMiddleware(logger)}
	for _, middleware := range middlewares {
		mux = middleware(mux)
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
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
