package server

import (
	"anomaly_detector/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	healthStatus = `{"status":"healthy"}`
)

type IHealthcheckServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type healthcheckServer struct {
	cfg    *config.InitConfig
	server *http.Server
}

// NewHealthcheckServer creates a lightweight HTTP server for health checks only
func NewHealthcheckServer(cfg *config.InitConfig) IHealthcheckServer {
	return &healthcheckServer{
		cfg:    cfg,
		server: NewHTTPHealthServer(cfg.HealthcheckPort),
	}
}

func NewHTTPHealthServer(port int) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: cReadHeaderTimeout,
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(healthStatus))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (hs *healthcheckServer) ListenAndServe() error {
	http.HandleFunc("/health", healthHandler)

	slog.Info("Healthcheck server is listening on", "port", hs.cfg.HealthcheckPort)

	if err := hs.server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("error while running healthcheck server: %w", err)
	}

	return nil
}

func (hs *healthcheckServer) Shutdown(ctx context.Context) error {
	if err := hs.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
