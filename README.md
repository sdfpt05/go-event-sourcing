# go-event-sourcing

A guide to learning event sourcing in Go by building an bank account system, following clean architecture principles.

## File Structure

.
├── cmd
│ └── main.go
├── domain
│ └── account.go
├── application
│ └── account_service.go
├── infrastructure
│ ├── postgres_event_store.go
│ ├── elasticsearch_projection.go
│ ├── azure_service_bus.go
│ └── multi_event_publisher.go
├── interfaces
│ ├── http_handler.go
│ └── event_consumer.go
├── go.mod
├── go.sum
└── README.md

## Components Overview

### Domain Layer (domain/account.go)

This layer represents the core business logic of the system.

- Defines the Account aggregate.
- Implements the events: `AccountCreatedEvent`, `DepositedEvent`, and `WithdrawnEvent`.
- Enforces business rules around deposits and withdrawals.

### Application Layer (application/account_service.go)

The Application layer contains the use cases that orchestrate the domain and infrastructure layers.

- `AccountService` handles account creation, deposits, and withdrawals.
- Defines interfaces for the `EventStore` and `EventPublisher`.

### Infrastructure Layer

This layer provides concrete implementations for external systems and databases.

- PostgreSQL event store (infrastructure/postgres_event_store.go)
- Elasticsearch projection for read models (infrastructure/elasticsearch_projection.go)
- Azure Service Bus publisher for sending events (infrastructure/azure_service_bus.go)
- Multi-event publisher that handles multiple event stores (infrastructure/multi_event_publisher.go)

### Interfaces Layer

Handles communication between the external world and the system.

- HTTP request handling via `http_handler.go`.
- Event consumption from Azure Service Bus via `event_consumer.go`.

### Main Application (cmd/main.go)

Entry point to the application.

- Initializes and configures all components.
- Starts the HTTP server.
- Sets up the event consumer for asynchronous processing.

## Setup and Running the Application

1. Install Go (version 1.16+ recommended).
2. Configure the required services:
   - PostgreSQL for the event store.
   - Elasticsearch for projections.
   - Azure Service Bus for event publishing.
3. Set up the following environment variables:
   - `POSTGRES_CONNECTION_STRING`
   - `ELASTICSEARCH_URL`
   - `SERVICEBUS_CONNECTION_STRING`
   - `SERVICEBUS_QUEUE_NAME`
4. Start the application:

```bash
   go run cmd/main.go
```

## API Endpoints

- `POST /account/create` - Creates a new account.
- `POST /account/deposit` - Deposits money into an account.
- `POST /account/withdraw` - Withdraws money from an account.
- `GET /account/balance?id={accountId}` - Retrieves the balance for a given account.

## Event Processing Workflow

Events such as `AccountCreated`, `Deposited`, and `Withdrawn` are published to Azure Service Bus. These events are processed asynchronously by consumers, allowing the system to scale and maintain resilience.

## Architecture Benefits

- **Separation of Concerns**: Cleanly separates domain, application, infrastructure, and interface logic.
- **Dependency Inversion**: Business logic is decoupled from external dependencies.
- **Event Sourcing**: Every state change is represented as an event.
- **Scalable Architecture**: Asynchronous event processing enables the system to scale effectively.
- **Testability**: The clean architecture makes unit testing and integration testing straightforward.

## Future Enhancements

[] Add authentication and authorization.
[] Implement logging and monitoring for better observability.
[] Introduce event versioning to handle changes in event structure.
[] Write unit and integration tests for each layer.
[] Implement snapshotting to improve performance for large event stores.
