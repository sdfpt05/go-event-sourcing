package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/sdfpt05/go-event-sourcing/internal/domain"
	"github.com/sdfpt05/go-event-sourcing/internal/repository"
)

type PostgresEventStore struct {
	queries *repository.Queries
}

func NewPostgresEventStore(queries *repository.Queries) *PostgresEventStore {
	return &PostgresEventStore{queries: queries}
}

func (s *PostgresEventStore) SaveEvent(event domain.Event) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.queries.SaveEvent(context.Background(), repository.SaveEventParams{
		AggregateID: event.AggregateID(),
		EventType:   event.EventType(),
		EventData:   eventData,
		Timestamp:   event.Timestamp(),
	})
}

func (s *PostgresEventStore) GetEvents(aggregateID string) ([]domain.Event, error) {
	events, err := s.queries.GetEvents(context.Background(), aggregateID)
	if err != nil {
		return nil, err
	}

	domainEvents := make([]domain.Event, len(events))
	for i, event := range events {
		var baseEvent domain.BaseEvent
		err := json.Unmarshal(event.EventData, &baseEvent)
		if err != nil {
			return nil, err
		}

		switch baseEvent.EventType() {
		case "AccountCreated":
			var e domain.AccountCreatedEvent
			err = json.Unmarshal(event.EventData, &e)
			domainEvents[i] = e
		case "Deposited":
			var e domain.DepositedEvent
			err = json.Unmarshal(event.EventData, &e)
			domainEvents[i] = e
		case "Withdrawn":
			var e domain.WithdrawnEvent
			err = json.Unmarshal(event.EventData, &e)
			domainEvents[i] = e
		}

		if err != nil {
			return nil, err
		}
	}

	return domainEvents, nil
}