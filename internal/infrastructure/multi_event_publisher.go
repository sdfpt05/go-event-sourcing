package infrastructure

import (
	"github.com/sdfpt05/go-event-sourcing/internal/application"
	"github.com/sdfpt05/go-event-sourcing/internal/domain"
)

type MultiEventPublisher struct {
	Publishers []application.EventPublisher
}

func (m *MultiEventPublisher) PublishEvent(event domain.Event) error {
	for _, publisher := range m.Publishers {
		if err := publisher.PublishEvent(event); err != nil {
			return err
		}
	}
	return nil
}
