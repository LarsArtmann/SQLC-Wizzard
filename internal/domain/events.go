package domain

import (
	"time"
)

// Event represents a domain event
type Event interface {
	// ID returns the unique identifier for this event
	ID() string

	// AggregateID returns the ID of the aggregate that generated this event
	AggregateID() string

	// EventType returns the type of this event
	EventType() string

	// OccurredAt returns when this event occurred
	OccurredAt() time.Time

	// Data returns the event data
	Data() any
}

// BaseEvent provides common functionality for all domain events
type BaseEvent struct {
	id          string
	aggregateID string
	eventType   string
	occurredAt  time.Time
	data        any
}

// NewBaseEvent creates a new base event
func NewBaseEvent(id, aggregateID, eventType string, data any) *BaseEvent {
	return &BaseEvent{
		id:          id,
		aggregateID: aggregateID,
		eventType:   eventType,
		occurredAt:  time.Now(),
		data:        data,
	}
}

// ID returns the unique identifier for this event
func (e *BaseEvent) ID() string {
	return e.id
}

// AggregateID returns the ID of the aggregate that generated this event
func (e *BaseEvent) AggregateID() string {
	return e.aggregateID
}

// EventType returns the type of this event
func (e *BaseEvent) EventType() string {
	return e.eventType
}

// OccurredAt returns when this event occurred
func (e *BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Data returns the event data
func (e *BaseEvent) Data() any {
	return e.data
}
