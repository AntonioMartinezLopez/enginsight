package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AntonioMartinezLopez/enginsight/jrpc"
	"github.com/AntonioMartinezLopez/enginsight/server/config"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/domain/counter"
	"github.com/AntonioMartinezLopez/enginsight/server/internal/infra/store"
)

func main() {
	// Load configuration from .env file
	config := config.LoadConfig()

	// set up application
	store := store.New()
	counterService := counter.New(store)

	// Create the JSON-RPC server with our service implementation
	handlers := jrpc.Handlers{
		CountService: counterService,
	}
	rpcServer := jrpc.NewServer(handlers)
	defer rpcServer.Close()

	// Set up HTTP server with the JSON-RPC handler
	mux := http.NewServeMux()
	mux.Handle(config.RPCPath, rpcServer)

	httpServer := &http.Server{
		Addr:              config.GetServerAddress(),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting JSON-RPC server on %s%s", config.GetServerAddress(), config.RPCPath)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
	count, err := counterService.GetNumberOfProcessedMessages(ctx)
	if err != nil {
		log.Fatalf("Failed to get number of processed messages: %v", err)
	}
	log.Printf("Total messages processed: %d", count)
}
