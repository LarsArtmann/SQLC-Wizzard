package events

import (
	"fmt"
	"maps"
)

// Event represents a domain event with typed data
// Replaces interface{} Data() method with typed EventData
type Event interface {
	ID() string
	Type() string
	AggregateID() string
	AggregateType() string
	Data() EventData
	OccurredAt() int64
	Version() int
	Metadata() map[string]string
}

// BaseEvent implements Event interface with typed event data
type BaseEvent struct {
	id            string
	typ           string
	aggregateID   string
	aggregateType string
	data          EventData
	occurredAt    int64
	version       int
	metadata      map[string]string
}

// NewBaseEvent creates a new base event with typed data and validation
func NewBaseEvent(id, eventType, aggregateID, aggregateType string, data EventData, version int) (*BaseEvent, error) {
	if id == "" {
		return nil, &EventValidationError{
			Code:    "EMPTY_EVENT_ID",
			Message: "Event ID cannot be empty",
		}
	}

	if eventType == "" {
		return nil, &EventValidationError{
			Code:    "EMPTY_EVENT_TYPE",
			Message: "Event type cannot be empty",
		}
	}

	if aggregateID == "" {
		return nil, &EventValidationError{
			Code:    "EMPTY_AGGREGATE_ID",
			Message: "Aggregate ID cannot be empty",
		}
	}

	if aggregateType == "" {
		return nil, &EventValidationError{
			Code:    "EMPTY_AGGREGATE_TYPE",
			Message: "Aggregate type cannot be empty",
		}
	}

	if data == nil {
		return nil, &EventValidationError{
			Code:    "NIL_EVENT_DATA",
			Message: "Event data cannot be nil",
		}
	}

	if version <= 0 {
		return nil, &EventValidationError{
			Code:    "INVALID_VERSION",
			Message: "Event version must be greater than 0",
		}
	}

	// Validate the event data
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("event data validation failed: %w", err)
	}

	return &BaseEvent{
		id:            id,
		typ:           eventType,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		data:          data,
		occurredAt:    data.OccurredAt().Unix(),
		version:       version,
		metadata:      make(map[string]string),
	}, nil
}

// ID returns the unique event identifier
func (be *BaseEvent) ID() string {
	return be.id
}

// Type returns the event type
func (be *BaseEvent) Type() string {
	return be.typ
}

// AggregateID returns the aggregate ID
func (be *BaseEvent) AggregateID() string {
	return be.aggregateID
}

// AggregateType returns the aggregate type
func (be *BaseEvent) AggregateType() string {
	return be.aggregateType
}

// Data returns the typed event data
func (be *BaseEvent) Data() EventData {
	return be.data
}

// OccurredAt returns the timestamp when the event occurred
func (be *BaseEvent) OccurredAt() int64 {
	return be.occurredAt
}

// Version returns the event version
func (be *BaseEvent) Version() int {
	return be.version
}

// Metadata returns event metadata
func (be *BaseEvent) Metadata() map[string]string {
	// Return a copy to prevent mutation
	copy := make(map[string]string)
	maps.Copy(copy, be.metadata)
	return copy
}

// WithMetadata adds metadata to the event
func (be *BaseEvent) WithMetadata(key, value string) *BaseEvent {
	be.metadata[key] = value
	return be
}

// WithMetadataBatch adds multiple metadata key-value pairs
func (be *BaseEvent) WithMetadataBatch(metadata map[string]string) *BaseEvent {
	maps.Copy(be.metadata, metadata)
	return be
}

// IsValid validates the event
func (be *BaseEvent) IsValid() bool {
	if be.id == "" {
		return false
	}
	if be.typ == "" {
		return false
	}
	if be.aggregateID == "" {
		return false
	}
	if be.aggregateType == "" {
		return false
	}
	if be.data == nil {
		return false
	}
	if be.version <= 0 {
		return false
	}

	// Validate the data
	if err := be.data.Validate(); err != nil {
		return false
	}

	return true
}
