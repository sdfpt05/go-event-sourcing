package infrastructure

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sdfpt05/go-event-sourcing/internal/domain"
	"time"
)

type PostgresEventStore struct {
	db *sql.DB
}

func NewPostgresEventStore(connectionString string) (*PostgresEventStore, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &PostgresEventStore{db: db}, nil
}

func (s *PostgresEventStore) SaveEvent(event domain.Event) error {
	dataJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		"INSERT INTO events (aggregate_id, type, data, timestamp) VALUES ($1, $2, $3, $4)",
		event.AggregateID(), event.EventType(), dataJSON, event.Timestamp(),
	)
	return err
}

func (s *PostgresEventStore) GetEvents(aggregateID string) ([]domain.Event, error) {
	rows, err := s.db.Query(
		"SELECT type, data, timestamp FROM events WHERE aggregate_id = $1 ORDER BY timestamp ASC",
		aggregateID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var eventType string
		var dataJSON []byte
		var timestamp time.Time
		err := rows.Scan(&eventType, &dataJSON, &timestamp)
		if err != nil {
			return nil, err
		}
		event, err := deserializeEvent(eventType, dataJSON, timestamp, aggregateID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func deserializeEvent(eventType string, dataJSON []byte, timestamp time.Time, aggregateID string) (domain.Event, error) {
	baseEvent := domain.NewBaseEvent(aggregateID, eventType, timestamp)
	switch eventType {
	case "AccountCreated":
		var e domain.AccountCreatedEvent
		e.BaseEvent = baseEvent
		err := json.Unmarshal(dataJSON, &e)
		return e, err
	case "Deposited":
		var e domain.DepositedEvent
		e.BaseEvent = baseEvent
		err := json.Unmarshal(dataJSON, &e)
		return e, err
	case "Withdrawn":
		var e domain.WithdrawnEvent
		e.BaseEvent = baseEvent
		err := json.Unmarshal(dataJSON, &e)
		return e, err
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}
