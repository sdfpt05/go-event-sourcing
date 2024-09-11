package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-service-bus-go"
	"github.com/sdfpt05/go-event-sourcing/domain"
)

type AzureServiceBusPublisher struct {
	sender *servicebus.Sender
}

func NewAzureServiceBusPublisher(connectionString, queueName string) (*AzureServiceBusPublisher, error) {
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		return nil, err
	}
	q, err := ns.NewQueue(queueName)
	if err != nil {
		return nil, err
	}
	return &AzureServiceBusPublisher{sender: q.NewSender()}, nil
}

func (p *AzureServiceBusPublisher) PublishEvent(event domain.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.sender.Send(context.Background(), servicebus.NewMessage(body))
}
