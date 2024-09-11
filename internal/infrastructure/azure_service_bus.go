package infrastructure

import (
	"context"
	"github.com/Azure/azure-service-bus-go"
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

	// Create a sender with the required context and options
	sender, err := q.NewSender(context.Background())
	if err != nil {
		return nil, err
	}

	return &AzureServiceBusPublisher{sender: sender}, nil
}
