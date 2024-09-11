// application/account_service.go
package application

import (
	"github.com/sdfpt05/go-event-sourcing/domain"
	"time"
)

type EventStore interface {
	SaveEvent(event domain.Event) error
	GetEvents(aggregateID string) ([]domain.Event, error)
}

type EventPublisher interface {
	PublishEvent(event domain.Event) error
}

type AccountService struct {
	eventStore     EventStore
	eventPublisher EventPublisher
}

func NewAccountService(eventStore EventStore, eventPublisher EventPublisher) *AccountService {
	return &AccountService{
		eventStore:     eventStore,
		eventPublisher: eventPublisher,
	}
}

func (s *AccountService) CreateAccount(id string, initialBalance float64) error {
	event := domain.AccountCreatedEvent{
		BaseEvent: domain.BaseEvent{
			aggregateID: id,
			eventType:   "AccountCreated",
			timestamp:   time.Now(),
		},
		InitialBalance: initialBalance,
	}
	err := s.eventStore.SaveEvent(event)
	if err != nil {
		return err
	}
	return s.eventPublisher.PublishEvent(event)
}

func (s *AccountService) Deposit(id string, amount float64) error {
	account, err := s.loadAccount(id)
	if err != nil {
		return err
	}
	event, err := account.Deposit(amount)
	if err != nil {
		return err
	}
	err = s.eventStore.SaveEvent(event)
	if err != nil {
		return err
	}
	return s.eventPublisher.PublishEvent(event)
}

func (s *AccountService) Withdraw(id string, amount float64) error {
	account, err := s.loadAccount(id)
	if err != nil {
		return err
	}
	event, err := account.Withdraw(amount)
	if err != nil {
		return err
	}
	err = s.eventStore.SaveEvent(event)
	if err != nil {
		return err
	}
	return s.eventPublisher.PublishEvent(event)
}

func (s *AccountService) GetBalance(id string) (float64, error) {
	account, err := s.loadAccount(id)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func (s *AccountService) loadAccount(id string) (*domain.Account, error) {
	events, err := s.eventStore.GetEvents(id)
	if err != nil {
		return nil, err
	}
	account := &domain.Account{ID: id}
	for _, event := range events {
		account.Apply(event)
	}
	return account, nil
}
