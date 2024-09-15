package interfaces

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Azure/azure-service-bus-go"
	"github.com/sdfpt05/go-event-sourcing/internal/application"
	"github.com/sdfpt05/go-event-sourcing/internal/domain"
)

type AzureServiceBusConsumer struct {
	receiver       *servicebus.Receiver
	accountService *application.AccountService
}

func NewAzureServiceBusConsumer(connectionString, queueName string, accountService *application.AccountService) (*AzureServiceBusConsumer, error) {
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		return nil, err
	}
	q, err := ns.NewQueue(queueName)
	if err != nil {
		return nil, err
	}

	receiver, err := q.NewReceiver(context.Background())
	if err != nil {
		return nil, err
	}

	return &AzureServiceBusConsumer{
		receiver:       receiver,
		accountService: accountService,
	}, nil
}

func (c *AzureServiceBusConsumer) Start(ctx context.Context) error {
	for {
		if err := c.receiver.ReceiveOne(ctx, servicebus.HandlerFunc(c.handleMessage)); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("Error receiving message: %v", err)
		}
	}
}

func (c *AzureServiceBusConsumer) Stop(ctx context.Context) error {
	return c.receiver.Close(ctx)
}

func (c *AzureServiceBusConsumer) handleMessage(ctx context.Context, msg *servicebus.Message) error {
	var event domain.Event
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return err
	}

	switch event.EventType() {
	case "AccountCreated":
		var e domain.AccountCreatedEvent
		if err := json.Unmarshal(msg.Data, &e); err != nil {
			return err
		}
		return c.accountService.CreateAccount(e.AggregateID(), e.InitialBalance)
	case "Deposited":
		var e domain.DepositedEvent
		if err := json.Unmarshal(msg.Data, &e); err != nil {
			return err
		}
		return c.accountService.Deposit(e.AggregateID(), e.Amount)
	case "Withdrawn":
		var e domain.WithdrawnEvent
		if err := json.Unmarshal(msg.Data, &e); err != nil {
			return err
		}
		return c.accountService.Withdraw(e.AggregateID(), e.Amount)
	default:
		log.Printf("Unknown event type: %s", event.EventType())
	}

	return msg.Complete(ctx)
}