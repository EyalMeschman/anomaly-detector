package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"anomaly_detector/config"
	"anomaly_detector/infrautils"

	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	container := buildContainer()

	err := container.Invoke(runServer)
	if err != nil {
		log.Panic(err)
	}
}

func buildContainer() *dig.Container {
	c := dig.New()

	// Register configuration
	infrautils.IocProvideWrapper(c, config.LoadInit)

	return c
}

func runServer(cfg *config.InitConfig, router *mux.Router) error {
	ctx := context.Background()

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	slog.InfoContext(ctx, "Starting HTTP server", "address", addr)

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	signals := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	slog.InfoContext(ctx, "Server is ready to accept connections", "address", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "Could not listen on address", "address", addr, "error", err)
		}
	}()

	go func() {
		sig := <-signals
		slog.InfoContext(ctx, "Signal received", "signal", sig)

		shutdown <- true
	}()

	<-shutdown

	doShutdown(server)

	slog.InfoContext(ctx, "Server exited")

	return nil
}

func doShutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.InfoContext(ctx, "shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Error during server shutdown", "error", err)
	}
}
