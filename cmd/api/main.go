// Package main provides the entry point for the KV store HTTP API server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kv-store-wal/internal/server"
	"kv-store-wal/internal/service"
)

const (
	defaultPort = "8080"
)

func main() {
	//  initialize in-memory store and wal connection
	err := service.Initialize()
	if err != nil {
		fmt.Println("Initialize failed, err: %w", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize server
	srv := server.NewHandler()

	// Set up HTTP server
	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second, // Prevent Slowloris attacks
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server starting on port %s...\n", port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Server shutting down...")

	fmt.Println("Closing wal...")
	err = service.CloseWAL()
	if err != nil {
		fmt.Println("error closing wal, err: ", err.Error())
	} else {
		fmt.Println("Wal closed successfully")

	}

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited properly")
}
