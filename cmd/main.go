// cmd/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sdfpt05/go-event-sourcing/application"
	"github.com/sdfpt05/go-event-sourcing/infrastructure"
	"github.com/sdfpt05/go-event-sourcing/interfaces"
)

func main() {
	// Initialize PostgreSQL event store
	postgresEventStore, err := infrastructure.NewPostgresEventStore(os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL event store: %v", err)
	}

	// Initialize Elasticsearch projection
	esProjection, err := infrastructure.NewElasticsearchProjection([]string{os.Getenv("ELASTICSEARCH_URL")})
	if err != nil {
		log.Fatalf("Failed to initialize Elasticsearch projection: %v", err)
	}

	// Initialize Azure Service Bus publisher
	azurePublisher, err := infrastructure.NewAzureServiceBusPublisher(
		os.Getenv("SERVICEBUS_CONNECTION_STRING"),
		os.Getenv("SERVICEBUS_QUEUE_NAME"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Azure Service Bus publisher: %v", err)
	}

	// Create a multi-publisher that publishes to both Elasticsearch and Azure Service Bus
	multiPublisher := &infrastructure.MultiEventPublisher{
		Publishers: []application.EventPublisher{esProjection, azurePublisher},
	}

	// Initialize account service
	accountService := application.NewAccountService(postgresEventStore, multiPublisher)

	// Initialize HTTP handler
	accountHandler := interfaces.NewAccountHandler(accountService)

	// Set up HTTP routes
	http.HandleFunc("/account/create", accountHandler.CreateAccount)
	http.HandleFunc("/account/deposit", accountHandler.Deposit)
	http.HandleFunc("/account/withdraw", accountHandler.Withdraw)
	http.HandleFunc("/account/balance", accountHandler.GetBalance)

	// Start HTTP server
	go func() {
		log.Printf("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Initialize and start event consumer
	eventConsumer, err := interfaces.NewEventConsumer(
		os.Getenv("SERVICEBUS_CONNECTION_STRING"),
		os.Getenv("SERVICEBUS_QUEUE_NAME"),
		accountService,
	)
	if err != nil {
		log.Fatalf("Failed to initialize event consumer: %v", err)
	}

	log.Println("Starting event consumer")
	eventConsumer.Start()
}
