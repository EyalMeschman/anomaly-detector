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

type HealthcheckServer struct {
	cfg    *config.InitConfig
	server *http.Server
}

// NewHealthCheckServer creates a lightweight HTTP server for health checks only
func NewHealthCheckServer(cfg *config.InitConfig) IHealthcheckServer {
	return &HealthcheckServer{
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

func (h *HealthcheckServer) ListenAndServe() error {
	http.HandleFunc("/health", healthHandler)

	slog.Info("Healthcheck server is listening on", "port", h.cfg.HealthcheckPort)

	if err := h.server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("error while running healthcheck server: %w", err)
	}

	return nil
}

func (h *HealthcheckServer) Shutdown(ctx context.Context) error {
	if err := h.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
