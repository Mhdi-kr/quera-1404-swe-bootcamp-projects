package domain

import "time"

// DomainEvent is a fact that happened inside a bounded context.
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

// EventRecorder helps aggregates track domain events until application services publish them.
type EventRecorder struct {
	events []DomainEvent
}

func (r *EventRecorder) Record(event DomainEvent) {
	if event == nil {
		return
	}
	r.events = append(r.events, event)
}

func (r *EventRecorder) PullEvents() []DomainEvent {
	pulled := make([]DomainEvent, len(r.events))
	copy(pulled, r.events)
	r.events = nil
	return pulled
}
