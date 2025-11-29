package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"anomaly_detector/config"
	"anomaly_detector/infrautils"
	"anomaly_detector/server"
	"anomaly_detector/store"
	"anomaly_detector/validator"

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

	// Register router
	infrautils.IocProvideWrapper(c, mux.NewRouter)

	// Register configuration
	infrautils.IocProvideWrapper(c, config.LoadInit)

	// Register server components
	infrautils.IocProvideWrapper(c, server.NewHTTPServer)
	infrautils.IocProvideWrapper(c, server.NewHealthCheckServer)

	// Register store
	infrautils.IocProvideWrapper(c, store.NewModelStore)

	// Register handlers
	infrautils.IocProvideWrapper(c, store.NewStoreHandler)
	infrautils.IocProvideWrapper(c, validator.NewValidateHandler)

	return c
}

func setMuxHandlers(
	router *mux.Router,
	storeHandler store.IStoreHandler,
	validateHandler validator.IValidateHandler,
) {
	router.HandleFunc("/models", storeHandler.Handle).Methods("POST")

	router.HandleFunc("/validate", validateHandler.Handle).Methods("POST")
}

func runServer(
	router *mux.Router, mainServer server.IHTTPServer, store store.IStoreHandler,
	validate validator.IValidateHandler, healthServer server.IHealthcheckServer) error {
	ctx := context.Background()

	signals := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)

	setMuxHandlers(router, store, validate)

	mainServer.SetHandler(router)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start health check server
	go func() {
		if err := healthServer.ListenAndServe(); err != nil {
			slog.ErrorContext(ctx, "Health check server error", "error", err)
		}
	}()

	// Start main server
	go func() {
		if err := mainServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "Main server error", "error", err)
		}
	}()

	slog.InfoContext(ctx, "Servers are up and running")

	// Wait for shutdown signal
	go func() {
		sig := <-signals
		slog.InfoContext(ctx, "Signal received", "signal", sig)

		shutdown <- true
	}()

	<-shutdown

	// Shutdown servers gracefully
	doShutdown(mainServer, healthServer)

	slog.InfoContext(ctx, "Servers exited")

	return nil
}

func doShutdown(mainServer server.IHTTPServer, healthServer server.IHealthcheckServer) {
	ctx := context.Background()
	slog.InfoContext(ctx, "Shutting down servers...")

	if err := mainServer.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Error during main server shutdown", "error", err)
	}

	if err := healthServer.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Error during healthcheck server shutdown", "error", err)
	}
}
