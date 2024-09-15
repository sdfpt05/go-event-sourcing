package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sdfpt05/go-event-sourcing/internal/domain"
)

type ElasticsearchPublisher struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearchPublisher(addresses []string, index string) (*ElasticsearchPublisher, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticsearchPublisher{
		client: client,
		index:  index,
	}, nil
}

func (p *ElasticsearchPublisher) PublishEvent(event domain.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.client.Index(
		p.index,
		&body,
		p.client.Index.WithDocumentID(event.AggregateID()),
		p.client.Index.WithContext(context.Background()),
	)

	return err
}