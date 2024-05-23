package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Rezab98/web-analyzer/internal/analyzer"
	"github.com/Rezab98/web-analyzer/internal/server"
	"github.com/Rezab98/web-analyzer/pkg/pagedownloader"
)

// The main function is the entry point of the application.
// It calls the run function and logs any errors encountered.
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// The run function sets up the application, starts the HTTP server, and handles graceful shutdown.
func run() error {

	// config
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("load config failed: %v", err)
	}

	// logger
	if err := configureLogger(&cfg.Logger); err != nil {
		return fmt.Errorf("can not configure logger: %v", err)
	}

	pageDownloader := pagedownloader.New()

	pageAnalyzerService := analyzer.New(pageDownloader)

	httpServer := server.New(
		&server.Config{
			Port: cfg.HTTPServer.Port,
			Host: cfg.HTTPServer.Host,
		},
		pageAnalyzerService,
	)

	// Create a context that will be canceled on shutdown signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the server in a goroutine
	go func() {
		logrus.Infof("Starting HTTP server on port %d...", cfg.HTTPServer.Port)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for the shutdown signal
	<-ctx.Done()
	logrus.Infof("Shutting down HTTP server gracefully...")

	// Perform graceful shutdown with a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logrus.Fatalf("HTTP shutdown error: %v", err)
	}
	logrus.Infof("Graceful shutdown complete.")

	return nil
}
