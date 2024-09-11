package domain

import (
	"testing"
	"time"
)

func TestAccount_Apply_AccountCreatedEvent(t *testing.T) {
	account := &Account{}
	event := AccountCreatedEvent{
		BaseEvent:      NewBaseEvent("123", "AccountCreated", time.Now()),
		InitialBalance: 100,
	}

	account.Apply(event)

	if account.ID != "123" {
		t.Errorf("Expected account ID to be '123', got '%s'", account.ID)
	}
	if account.Balance != 100 {
		t.Errorf("Expected account balance to be 100, got %f", account.Balance)
	}
}

func TestAccount_Apply_DepositedEvent(t *testing.T) {
	account := &Account{ID: "123", Balance: 100}
	event := DepositedEvent{
		BaseEvent: NewBaseEvent("123", "Deposited", time.Now()),
		Amount:    50,
	}

	account.Apply(event)

	if account.Balance != 150 {
		t.Errorf("Expected account balance to be 150, got %f", account.Balance)
	}
}

func TestAccount_Apply_WithdrawnEvent(t *testing.T) {
	account := &Account{ID: "123", Balance: 100}
	event := WithdrawnEvent{
		BaseEvent: NewBaseEvent("123", "Withdrawn", time.Now()),
		Amount:    30,
	}

	account.Apply(event)

	if account.Balance != 70 {
		t.Errorf("Expected account balance to be 70, got %f", account.Balance)
	}
}

func TestAccount_Deposit(t *testing.T) {
	account := &Account{ID: "123", Balance: 100}

	event, err := account.Deposit(50)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	depositedEvent, ok := event.(DepositedEvent)
	if !ok {
		t.Errorf("Expected DepositedEvent, got %T", event)
	}

	if depositedEvent.Amount != 50 {
		t.Errorf("Expected deposit amount to be 50, got %f", depositedEvent.Amount)
	}

	if account.Balance != 150 {
		t.Errorf("Expected account balance to be 150, got %f", account.Balance)
	}
}

func TestAccount_Deposit_NegativeAmount(t *testing.T) {
	account := &Account{ID: "123", Balance: 100}

	_, err := account.Deposit(-50)
	if err == nil {
		t.Error("Expected error for negative deposit amount, got nil")
	}
}

func TestAccount_Withdraw(t *testing.T) {
	account := &Account{ID: "123", Balance: 100}

	event, err := account.Withdraw(30)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	withdrawnEvent, ok := event.(WithdrawnEvent)
	if !ok {
		t.Errorf("Expected WithdrawnEvent, got %T", event)
	}

	if withdrawnEvent.Amount != 30 {
		t.Errorf("Expected withdrawal amount to be 30, got %f", withdrawnEvent.Amount)
	}

	if account.Balance != 70 {
		t.Errorf("Expected account balance to be 70, got %f", account.Balance)
	}
}

func TestAccount_Withdraw_InsufficientFunds(t *testing.T) {
	account := &Account{ID: "123", Balance: 100}

	_, err := account.Withdraw(150)
	if err == nil {
		t.Error("Expected error for insufficient funds, got nil")
	}
}

func TestAccount_Withdraw_NegativeAmount(t *testing.T) {
	account := &Account{ID: "123", Balance: 100}

	_, err := account.Withdraw(-50)
	if err == nil {
		t.Error("Expected error for negative withdrawal amount, got nil")
	}
}

func TestCreateAccount(t *testing.T) {
	account, event := CreateAccount("123", 100)

	if account.ID != "123" {
		t.Errorf("Expected account ID to be '123', got '%s'", account.ID)
	}
	if account.Balance != 100 {
		t.Errorf("Expected account balance to be 100, got %f", account.Balance)
	}

	createdEvent, ok := event.(AccountCreatedEvent)
	if !ok {
		t.Errorf("Expected AccountCreatedEvent, got %T", event)
	}

	if createdEvent.AggregateID() != "123" {
		t.Errorf("Expected event aggregate ID to be '123', got '%s'", createdEvent.AggregateID())
	}
	if createdEvent.InitialBalance != 100 {
		t.Errorf("Expected initial balance to be 100, got %f", createdEvent.InitialBalance)
	}
}
