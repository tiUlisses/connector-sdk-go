package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server exposes optional operational endpoints.
type Server struct {
	httpServer *http.Server
}

// New returns an embedded HTTP server with /healthz and /metrics handlers.
func New(addr string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/metrics", promhttp.Handler())

	return &Server{httpServer: &http.Server{Addr: addr, Handler: mux}}
}

// Start begins serving in a goroutine.
func (s *Server) Start() error {
	go func() {
		_ = s.httpServer.ListenAndServe()
	}()
	return nil
}

// Stop gracefully stops the server.
func (s *Server) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	return nil
}
