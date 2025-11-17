package events

import (
	"maps"
	"time"
)

// EventData represents typed event data interface
// Replaces 'any' type with proper type constraints
type EventData interface {
	EventType() string
	OccurredAt() time.Time
	AggregateID() string
	Data() map[string]any // Safe map for JSON serialization
	Validate() error
}

// BaseEventData provides common event data functionality
type BaseEventData struct {
	eventType   string
	occurredAt  time.Time
	aggregateID string
	data        map[string]any
}

// NewBaseEventData creates a new base event data with validation
func NewBaseEventData(eventType, aggregateID string, data map[string]any) (*BaseEventData, error) {
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

	if data == nil {
		data = make(map[string]any)
	}

	return &BaseEventData{
		eventType:   eventType,
		occurredAt:  time.Now(),
		aggregateID: aggregateID,
		data:        data,
	}, nil
}

// EventType returns the event type
func (bed *BaseEventData) EventType() string {
	return bed.eventType
}

// OccurredAt returns when the event occurred
func (bed *BaseEventData) OccurredAt() time.Time {
	return bed.occurredAt
}

// AggregateID returns the aggregate ID
func (bed *BaseEventData) AggregateID() string {
	return bed.aggregateID
}

// Data returns the event data as a safe map
func (bed *BaseEventData) Data() map[string]any {
	// Return a copy to prevent mutation
	copy := make(map[string]any)
	maps.Copy(copy, bed.data)
	return copy
}

// Validate validates the base event data
func (bed *BaseEventData) Validate() error {
	if bed.eventType == "" {
		return &EventValidationError{
			Code:    "INVALID_EVENT_TYPE",
			Message: "Event type is required",
		}
	}

	if bed.aggregateID == "" {
		return &EventValidationError{
			Code:    "INVALID_AGGREGATE_ID",
			Message: "Aggregate ID is required",
		}
	}

	if bed.occurredAt.IsZero() {
		return &EventValidationError{
			Code:    "INVALID_OCCURRED_AT",
			Message: "Occurred at time is required",
		}
	}

	return nil
}

// EventValidationError represents event validation errors
type EventValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *EventValidationError) Error() string {
	return e.Message
}
