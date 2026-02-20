package application

import (
	"context"

	shareddomain "richdomainmodeling/internal/shared/domain"
)

type EventHandler func(context.Context, shareddomain.DomainEvent) error

type EventBus interface {
	Subscribe(eventName string, handler EventHandler)
	Publish(ctx context.Context, events ...shareddomain.DomainEvent) error
}
