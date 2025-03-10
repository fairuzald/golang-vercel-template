package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-template/api"
	"golang-template/configs"
	"golang-template/infrastructure/firebase"
	"golang-template/infrastructure/logger"
)

// @title GolangTemplate API
func main() {
	// Setup configuration
	cfg := configs.LoadConfig()

	// Initialize logger
	log := logger.NewLogger(cfg.Debug)
	log.Info("Starting API server", "name", cfg.AppName, "env", cfg.Environment)

	// Initialize Firebase client
	fbClient, err := firebase.Initialize(cfg, log)
	if err != nil {
		log.Warn("Firebase initialization failed", "error", err)
	} else {
		defer fbClient.Close()
	}

	// Setup HTTP router
	router := api.SetupRouter(cfg, log)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Info("Server listening", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server exited gracefully")
}
