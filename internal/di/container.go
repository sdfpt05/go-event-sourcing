package di

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sdfpt05/go-event-sourcing/internal/application"
	"github.com/sdfpt05/go-event-sourcing/internal/infrastructure"
	"github.com/sdfpt05/go-event-sourcing/internal/interfaces"
	"github.com/sdfpt05/go-event-sourcing/internal/repository"
)

type Container struct {
	DbPool         *pgxpool.Pool
	EventStore     application.EventStore
	AccountService *application.AccountService
	AccountHandler *interfaces.AccountHandler
}

func NewContainer(ctx context.Context, postgresURL string, esAddresses []string, esIndex string) (*Container, error) {
	dbPool, err := pgxpool.Connect(ctx, postgresURL)
	if err != nil {
		return nil, err
	}

	queries := repository.New(dbPool)
	eventStore := infrastructure.NewPostgresEventStore(queries)

	esPublisher, err := infrastructure.NewElasticsearchPublisher(esAddresses, esIndex)
	if err != nil {
		return nil, err
	}

	accountService := application.NewAccountService(eventStore, esPublisher)
	accountHandler := interfaces.NewAccountHandler(accountService)

	return &Container{
		DbPool:         dbPool,
		EventStore:     eventStore,
		AccountService: accountService,
		AccountHandler: accountHandler,
	}, nil
}

func (c *Container) Close() {
	if c.DbPool != nil {
		c.DbPool.Close()
	}
}