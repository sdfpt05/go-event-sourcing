package application

import (
	"errors"
	"testing"
	"time"

	"github.com/sdfpt05/go-event-sourcing/internal/domain"
)

// MockEventStore is a mock implementation of EventStore
type MockEventStore struct {
	SaveEventFunc func(event domain.Event) error
	GetEventsFunc func(aggregateID string) ([]domain.Event, error)
}

func (m *MockEventStore) SaveEvent(event domain.Event) error {
	return m.SaveEventFunc(event)
}

func (m *MockEventStore) GetEvents(aggregateID string) ([]domain.Event, error) {
	return m.GetEventsFunc(aggregateID)
}

// MockEventPublisher is a mock implementation of EventPublisher
type MockEventPublisher struct {
	PublishEventFunc func(event domain.Event) error
}

func (m *MockEventPublisher) PublishEvent(event domain.Event) error {
	return m.PublishEventFunc(event)
}

func TestAccountService_CreateAccount(t *testing.T) {
	mockEventStore := &MockEventStore{
		SaveEventFunc: func(event domain.Event) error {
			return nil
		},
	}
	mockEventPublisher := &MockEventPublisher{
		PublishEventFunc: func(event domain.Event) error {
			return nil
		},
	}

	service := NewAccountService(mockEventStore, mockEventPublisher)

	err := service.CreateAccount("123", 100)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAccountService_CreateAccount_SaveEventError(t *testing.T) {
	mockEventStore := &MockEventStore{
		SaveEventFunc: func(event domain.Event) error {
			return errors.New("save event error")
		},
	}
	mockEventPublisher := &MockEventPublisher{
		PublishEventFunc: func(event domain.Event) error {
			return nil
		},
	}

	service := NewAccountService(mockEventStore, mockEventPublisher)

	err := service.CreateAccount("123", 100)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAccountService_Deposit(t *testing.T) {
	mockEventStore := &MockEventStore{
		GetEventsFunc: func(aggregateID string) ([]domain.Event, error) {
			return []domain.Event{
				domain.AccountCreatedEvent{
					BaseEvent:      domain.NewBaseEvent("123", "AccountCreated", time.Now()),
					InitialBalance: 100,
				},
			}, nil
		},
		SaveEventFunc: func(event domain.Event) error {
			return nil
		},
	}
	mockEventPublisher := &MockEventPublisher{
		PublishEventFunc: func(event domain.Event) error {
			return nil
		},
	}

	service := NewAccountService(mockEventStore, mockEventPublisher)

	err := service.Deposit("123", 50)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAccountService_Withdraw(t *testing.T) {
	mockEventStore := &MockEventStore{
		GetEventsFunc: func(aggregateID string) ([]domain.Event, error) {
			return []domain.Event{
				domain.AccountCreatedEvent{
					BaseEvent:      domain.NewBaseEvent("123", "AccountCreated", time.Now()),
					InitialBalance: 100,
				},
			}, nil
		},
		SaveEventFunc: func(event domain.Event) error {
			return nil
		},
	}
	mockEventPublisher := &MockEventPublisher{
		PublishEventFunc: func(event domain.Event) error {
			return nil
		},
	}

	service := NewAccountService(mockEventStore, mockEventPublisher)

	err := service.Withdraw("123", 50)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAccountService_GetBalance(t *testing.T) {
	mockEventStore := &MockEventStore{
		GetEventsFunc: func(aggregateID string) ([]domain.Event, error) {
			return []domain.Event{
				domain.AccountCreatedEvent{
					BaseEvent:      domain.NewBaseEvent("123", "AccountCreated", time.Now()),
					InitialBalance: 100,
				},
				domain.DepositedEvent{
					BaseEvent: domain.NewBaseEvent("123", "Deposited", time.Now()),
					Amount:    50,
				},
			}, nil
		},
	}
	mockEventPublisher := &MockEventPublisher{}

	service := NewAccountService(mockEventStore, mockEventPublisher)

	balance, err := service.GetBalance("123")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if balance != 150 {
		t.Errorf("Expected balance to be 150, got %f", balance)
	}
}
