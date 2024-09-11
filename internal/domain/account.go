package domain

import (
	"errors"
	"time"
)

type Account struct {
	ID      string
	Balance float64
}

type Event interface {
	AggregateID() string
	EventType() string
	Timestamp() time.Time
}

type BaseEvent struct {
	aggregateID string
	eventType   string
	timestamp   time.Time
}

func (e BaseEvent) AggregateID() string  { return e.aggregateID }
func (e BaseEvent) EventType() string    { return e.eventType }
func (e BaseEvent) Timestamp() time.Time { return e.timestamp }

func NewBaseEvent(aggregateID, eventType string, timestamp time.Time) BaseEvent {
	return BaseEvent{
		aggregateID: aggregateID,
		eventType:   eventType,
		timestamp:   timestamp,
	}
}

type AccountCreatedEvent struct {
	BaseEvent
	InitialBalance float64
}

type DepositedEvent struct {
	BaseEvent
	Amount float64
}

type WithdrawnEvent struct {
	BaseEvent
	Amount float64
}

func (a *Account) Apply(event Event) {
	switch e := event.(type) {
	case AccountCreatedEvent:
		a.ID = e.AggregateID()
		a.Balance = e.InitialBalance
	case DepositedEvent:
		a.Balance += e.Amount
	case WithdrawnEvent:
		a.Balance -= e.Amount
	}
}

func (a *Account) Deposit(amount float64) (Event, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be positive")
	}
	event := DepositedEvent{
		BaseEvent: NewBaseEvent(a.ID, "Deposited", time.Now()),
		Amount:    amount,
	}
	a.Apply(event)
	return event, nil
}

func (a *Account) Withdraw(amount float64) (Event, error) {
	if amount <= 0 {
		return nil, errors.New("withdrawal amount must be positive")
	}
	if a.Balance < amount {
		return nil, errors.New("insufficient funds")
	}
	event := WithdrawnEvent{
		BaseEvent: NewBaseEvent(a.ID, "Withdrawn", time.Now()),
		Amount:    amount,
	}
	a.Apply(event)
	return event, nil
}

func CreateAccount(id string, initialBalance float64) (Account, Event) {
	account := Account{ID: id, Balance: initialBalance}
	event := AccountCreatedEvent{
		BaseEvent:      NewBaseEvent(id, "AccountCreated", time.Now()),
		InitialBalance: initialBalance,
	}
	account.Apply(event)
	return account, event
}
