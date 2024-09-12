package http

import (
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(next http.Handler) http.Handler

type loggerWriter struct {
	http.ResponseWriter
	code int
	size int
}

func (w *loggerWriter) Write(bytes []byte) (int, error) {
	n, err := w.ResponseWriter.Write(bytes)
	w.size += n
	return n, err
}

func (w *loggerWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LoggerMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := &loggerWriter{ResponseWriter: w}
			start := time.Now()
			next.ServeHTTP(lw, r)
			logger.Info(r.RemoteAddr,
				"method", r.Method,
				"url", r.URL.String(),
				"proto", r.Proto,
				"code", lw.code,
				"size", lw.size,
				"time", time.Since(start))
		})
	}
}

func RecovererMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rc := recover(); rc != nil {
					w.WriteHeader(http.StatusInternalServerError)
					if err, ok := rc.(error); ok {
						logger.Error(r.RemoteAddr, "err", err)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
