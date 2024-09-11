# go-event-sourcing

Learning event sourcing by building an account system using clean architecture principles.

## File Structure

```
.
├── cmd
│   └── main.go
├── domain
│   └── account.go
├── application
│   └── account_service.go
├── infrastructure
│   ├── postgres_event_store.go
│   ├── elasticsearch_projection.go
│   ├── azure_service_bus.go
│   └── multi_event_publisher.go
├── interfaces
│   ├── http_handler.go
│   └── event_consumer.go
├── go.mod
├── go.sum
└── README.md
```

## Components

### Domain Layer (domain/account.go)

Contains the core business logic, including:

- Account aggregate
- Event definitions (AccountCreatedEvent, DepositedEvent, WithdrawnEvent)
- Business rules for deposits and withdrawals

### Application Layer (application/account_service.go)

Implements use cases and coordinates between domain and infrastructure:

- AccountService for creating accounts, making deposits, and withdrawals
- Interfaces for EventStore and EventPublisher

### Infrastructure Layer

Provides concrete implementations for external services:

- PostgreSQL event store (infrastructure/postgres_event_store.go)
- Elasticsearch projection (infrastructure/elasticsearch_projection.go)
- Azure Service Bus publisher (infrastructure/azure_service_bus.go)
- Multi-event publisher (infrastructure/multi_event_publisher.go)

### Interfaces Layer

Handles external interactions:

- HTTP request handling (interfaces/http_handler.go)
- Azure Service Bus message consumption (interfaces/event_consumer.go)

### Main Application (cmd/main.go)

Ties everything together:

- Initializes all components
- Sets up HTTP server
- Starts event consumer

## Setup and Running

1. Ensure you have Go installed (version 1.16+ recommended).
2. Set up PostgreSQL, Elasticsearch, and Azure Service Bus.
3. Set the following environment variables:
   - POSTGRES_CONNECTION_STRING
   - ELASTICSEARCH_URL
   - SERVICEBUS_CONNECTION_STRING
   - SERVICEBUS_QUEUE_NAME
4. Run the application:
   ```
   go run cmd/main.go
   ```

## API Endpoints

- POST /account/create - Create a new account
- POST /account/deposit - Make a deposit
- POST /account/withdraw - Make a withdrawal
- GET /account/balance?id={accountId} - Get account balance

## Event Processing

Events are published to Azure Service Bus and consumed asynchronously. This allows for scalable and resilient event processing.

## Architecture Benefits

- Separation of Concerns: Each layer has a specific responsibility.
- Dependency Inversion: Core business logic doesn't depend on external frameworks.
- Event Sourcing: All state changes are represented as events.
- Message-Driven: Asynchronous event processing for scalability.
- Clean Architecture: Easy to test, maintain, and evolve.

## Future Improvements

[] Implement authentication and authorization
[] Add comprehensive logging and monitoring
[] Implement event versioning for schema evolution
[] Add unit and integration tests
[] Implement snapshotting for performance optimization
