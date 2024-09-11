package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sdfpt05/go-event-sourcing/internal/application"
	"github.com/sdfpt05/go-event-sourcing/internal/domain"
	"github.com/sdfpt05/go-event-sourcing/internal/infrastructure"
	"github.com/sdfpt05/go-event-sourcing/internal/interfaces"
)

// MockEventPublisher is a mock implementation of EventPublisher
type MockEventPublisher struct {
	PublishEventFunc func(event domain.Event) error
}

func (m *MockEventPublisher) PublishEvent(event domain.Event) error {
	return m.PublishEventFunc(event)
}

func TestIntegration_CreateAccountAndDeposit(t *testing.T) {
	// Set up a test database
	dbConnectionString := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	eventStore, err := infrastructure.NewPostgresEventStore(dbConnectionString)
	if err != nil {
		t.Fatalf("Failed to create event store: %v", err)
	}

	// Create a mock event publisher
	mockPublisher := &MockEventPublisher{
		PublishEventFunc: func(event domain.Event) error {
			return nil
		},
	}

	// Set up the account service
	accountService := application.NewAccountService(eventStore, mockPublisher)

	// Set up the HTTP handler
	handler := interfaces.NewAccountHandler(accountService)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(handler.CreateAccount))
	defer server.Close()

	// Test creating an account
	createAccountPayload := []byte(`{"id": "123", "initial_balance": 100}`)
	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createAccountPayload))
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", resp.Status)
	}

	// Set up a new test server for deposit
	depositServer := httptest.NewServer(http.HandlerFunc(handler.Deposit))
	defer depositServer.Close()

	// Test making a deposit
	depositPayload := []byte(`{"id": "123", "amount": 50}`)
	resp, err = http.Post(depositServer.URL, "application/json", bytes.NewBuffer(depositPayload))
	if err != nil {
		t.Fatalf("Failed to make deposit: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Verify the balance
	balanceServer := httptest.NewServer(http.HandlerFunc(handler.GetBalance))
	defer balanceServer.Close()

	resp, err = http.Get(balanceServer.URL + "?id=123")
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var balanceResponse struct {
		Balance float64 `json:"balance"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&balanceResponse); err != nil {
		t.Fatalf("Failed to decode balance response: %v", err)
	}

	if balanceResponse.Balance != 150 {
		t.Errorf("Expected balance to be 150, got %f", balanceResponse.Balance)
	}
}
