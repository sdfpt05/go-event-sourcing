# go-event-sourcing

A guide to learning event sourcing in Go by building an bank account system, following clean architecture principles.

## Features

- Event sourcing architecture
- PostgreSQL for event storage
- Elasticsearch for event querying
- Azure Service Bus for event consumption
- RESTful API for account operations

## Architecture Benefits

- **Separation of Concerns**: Cleanly separates domain, application, infrastructure, and interface logic.
- **Dependency Inversion**: Business logic is decoupled from external dependencies.
- **Event Sourcing**: Every state change is represented as an event.
- **Scalable Architecture**: Asynchronous event processing enables the system to scale effectively.
- **Testability**: The clean architecture makes unit testing and integration testing straightforward.

## Prerequisites

- Go 1.16+
- PostgreSQL
- Elasticsearch
- Azure Service Bus

## Project Structure

```
.
├── bin/
├── cmd/
│   └── main.go
├── internal/
│   ├── application/
│   ├── config/
│   ├── di/
│   ├── domain/
│   ├── infrastructure/
│   ├── interfaces/
│   ├── middleware/
│   └── repository/
├── tests/
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── sqlc.yaml
```

## API Endpoints

- `POST /account/create` - Creates a new account.
- `POST /account/deposit` - Deposits money into an account.
- `POST /account/withdraw` - Withdraws money from an account.
- `GET /account/balance?id={accountId}` - Retrieves the balance for a given account.

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go-event-sourcing.git
   cd go-event-sourcing
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set up environment variables:
   ```
   export POSTGRES_URL="postgres://user:password@localhost:5432/dbname"
   export ELASTICSEARCH_URL="http://localhost:9200"
   export SERVICEBUS_CONNECTION_STRING="your_connection_string"
   export SERVICEBUS_QUEUE_NAME="your_queue_name"
   export SERVER_PORT="8080"
   ```

4. Generate SQL code:
   ```
   sqlc generate
   ```

5. Build the application:
   ```
   go build -o bin/go-event-sourcing cmd/main.go
   ```

## Running the Application

Run the compiled binary:

```
./bin/go-event-sourcing
```

The application will start an HTTP server on the specified port and begin consuming events from Azure Service Bus.

## API Endpoints

- POST /api/v1/account/create - Create a new account
- POST /api/v1/account/deposit - Deposit funds
- POST /api/v1/account/withdraw - Withdraw funds
- GET /api/v1/account/balance - Get account balance

## Docker

To build and run the application using Docker, see the Dockerfile in the project root.

## Testing

Run the tests with:

```
go test ./...
```

## Future Enhancements


- [x] RESTful API for account operations
- [x] SQL query generation using sqlc
- [x] Middleware for logging and error handling
- [ ] Introduce event versioning to handle changes in event structure
- [ ] Write unit and integration tests for each layer
- [ ] Implement Snapshotting Improve performance for large event stores
- [ ] Implement logging and monitoring for better observability
- [ ] Implement event upcasting for handling multiple event versions
- [ ] Add support for event replay and system recovery
- [ ] Implement a projection builder for creating different views of the data
- [ ] Add support for event scheduling and time-based operations
- [ ] Implement a saga pattern for managing distributed transactions
- [ ] Create a dashboard for visualizing system metrics and event flows


## License

This project is licensed under the MIT License - see the LICENSE file for details.