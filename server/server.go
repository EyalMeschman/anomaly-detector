package server

import (
	"anomaly_detector/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const cReadHeaderTimeout = 5 * time.Second

type IHTTPServer interface {
	ListenAndServe() error
	SetHandler(handler http.Handler)
	Shutdown(ctx context.Context) error
}

type httpServer struct {
	srv  *http.Server
	port int
}

func NewHTTPServer(cfg *config.InitConfig) IHTTPServer {
	return &httpServer{
		srv:  &http.Server{ReadHeaderTimeout: cReadHeaderTimeout},
		port: cfg.ServerPort,
	}
}

func (s *httpServer) ListenAndServe() error {
	s.srv.Addr = fmt.Sprintf(":%d", s.port)

	slog.Info("HTTP server is listening on", "port", s.port)

	e := s.srv.ListenAndServe()

	if e == http.ErrServerClosed {
		return nil
	}

	return e
}

func (s *httpServer) SetHandler(handler http.Handler) {
	s.srv.Handler = handler
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
