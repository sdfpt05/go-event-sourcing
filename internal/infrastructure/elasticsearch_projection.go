package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sdfpt05/go-event-sourcing/internal/domain"
)

type ElasticsearchProjection struct {
	client *elasticsearch.Client
}

func NewElasticsearchProjection(addresses []string) (*ElasticsearchProjection, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: addresses,
	})
	if err != nil {
		return nil, err
	}
	return &ElasticsearchProjection{client: client}, nil
}

func (p *ElasticsearchProjection) PublishEvent(event domain.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = p.client.Index(
		"account_events",
		strings.NewReader(string(body)),
		p.client.Index.WithRefresh("true"),
	)
	return err
}
