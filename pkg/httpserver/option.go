package httpserver

import (
	"time"
)

type Option func(*Server)

func Addr(addr string) Option {
	return func(s *Server) {
		s.server.Addr = addr
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}
