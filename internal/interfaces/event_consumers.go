package interfaces

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-service-bus-go"
	"github.com/sdfpt05/go-event-sourcing/application"
	"github.com/sdfpt05/go-event-sourcing/domain"
	"log"
)

type EventConsumer struct {
	receiver       *servicebus.Receiver
	accountService *application.AccountService
}

func NewEventConsumer(connectionString, queueName string, accountService *application.AccountService) (*EventConsumer, error) {
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		return nil, err
	}
	q, err := ns.NewQueue(queueName)
	if err != nil {
		return nil, err
	}
	return &EventConsumer{
		receiver:       q.NewReceiver(),
		accountService: accountService,
	}, nil
}

func (c *EventConsumer) Start() {
	for {
		if err := c.receiver.ReceiveOne(context.Background(), servicebus.HandlerFunc(c.handleMessage)); err != nil {
			log.Printf("Error receiving message: %v", err)
		}
	}
}

func (c *EventConsumer) handleMessage(ctx context.Context, msg *servicebus.Message) error {
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
