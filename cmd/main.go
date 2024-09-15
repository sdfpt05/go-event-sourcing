package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sdfpt05/go-event-sourcing/internal/config"
	"github.com/sdfpt05/go-event-sourcing/internal/di"
	"github.com/sdfpt05/go-event-sourcing/internal/interfaces"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	container, err := di.NewContainer(
		context.Background(), 
		cfg.PostgresURL, 
		[]string{cfg.ElasticsearchURL}, 
		"events"
	)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	router := interfaces.NewRouter(container.AccountHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Starting HTTP server on :%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Initialize and start Azure Service Bus consumer
	consumer, err := interfaces.NewAzureServiceBusConsumer(
		cfg.ServiceBusConnStr,
		cfg.ServiceBusQueue,
		container.AccountService,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Azure Service Bus consumer: %v", err)
	}

	go func() {
		log.Println("Starting Azure Service Bus consumer")
		if err := consumer.Start(context.Background()); err != nil {
			log.Fatalf("Azure Service Bus consumer error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := consumer.Stop(ctx); err != nil {
		log.Fatalf("Failed to stop Azure Service Bus consumer: %v", err)
	}

	log.Println("Server exiting")
}
