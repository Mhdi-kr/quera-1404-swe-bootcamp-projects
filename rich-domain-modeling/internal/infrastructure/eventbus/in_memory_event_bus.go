package eventbus

import (
	"context"
	"fmt"
	"sync"

	sharedapp "richdomainmodeling/internal/shared/application"
	shareddomain "richdomainmodeling/internal/shared/domain"
)

type InMemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]sharedapp.EventHandler
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{subscribers: make(map[string][]sharedapp.EventHandler)}
}

func (b *InMemoryEventBus) Subscribe(eventName string, handler sharedapp.EventHandler) {
	if handler == nil || eventName == "" {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[eventName] = append(b.subscribers[eventName], handler)
}

func (b *InMemoryEventBus) Publish(ctx context.Context, events ...shareddomain.DomainEvent) error {
	for _, event := range events {
		if event == nil {
			continue
		}

		b.mu.RLock()
		handlers := append([]sharedapp.EventHandler(nil), b.subscribers[event.EventName()]...)
		b.mu.RUnlock()

		for _, handler := range handlers {
			if err := handler(ctx, event); err != nil {
				return fmt.Errorf("handling %q event: %w", event.EventName(), err)
			}
		}
	}

	return nil
}
