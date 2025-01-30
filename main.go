package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/putdotio/pas/internal/analytics"
	"github.com/putdotio/pas/internal/config"
	"github.com/putdotio/pas/internal/handler"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
)

var (
	version = "v0.0.0"
	commit  = "none"
	date    = "unknown"
)

var (
	versionFlag = flag.Bool("version", false, "version")
	configPath  = flag.String("config", "config.toml", "config file path")
)

func main() {
	fmt.Printf("Starting PAS version: %s, commit: %s, built: %s", version, commit, date)

	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		return
	}

	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", config.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// main server
	analytics := analytics.New(db, config.Secret, config.User, config.Events)
	handler := handler.New(analytics)
	server := http.Server{
		Addr:    config.ListenAddress,
		Handler: cors.Default().Handler(handler),
	}

	// metrics server
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := http.Server{
		Addr:    config.ListenAddressForMetrics,
		Handler: metricsMux,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create root errgroup with cancellable context
	g, ctx := errgroup.WithContext(ctx)

	// Run main server
	g.Go(func() error {
		log.Printf("Starting HTTP server on %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("main server error: %w", err)
		}
		return nil
	})

	// Run metrics server
	g.Go(func() error {
		log.Printf("Starting metrics server on %v", metricsServer.Addr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("metrics server error: %w", err)
		}
		return nil
	})

	// Handle shutdown signal
	g.Go(func() error {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		select {
		case <-stop:
			log.Println("Shutdown signal received")
			cancel() // Cancel context to trigger shutdown
		case <-ctx.Done():
			// Context was cancelled by server error
		}
		return nil
	})

	// Handle context cancellation and server shutdown
	g.Go(func() error {
		<-ctx.Done() // Wait for cancellation (either from signal or error)

		shutdownTimeout := time.Duration(config.ShutdownTimeout) * time.Millisecond

		// Shutdown main server with its own timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		log.Printf("Shutting down HTTP server with timeout: %v", shutdownTimeout)
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Main server shutdown error: %v", err)
		}
		cancel()

		// Shutdown metrics server with fresh timeout
		shutdownCtx, cancel = context.WithTimeout(context.Background(), shutdownTimeout)
		log.Printf("Shutting down metrics server with timeout: %v", shutdownTimeout)
		if err := metricsServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Metrics server shutdown error: %v", err)
		}
		cancel()

		return nil
	})

	// Wait for all goroutines to complete or for an error
	if err := g.Wait(); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}

	log.Println("Servers stopped")
}
