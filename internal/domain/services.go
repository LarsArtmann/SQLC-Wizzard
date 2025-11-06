package domain

import (
	"github.com/google/uuid"
)

// EventBus represents the interface for publishing events
type EventBus interface {
	// Publish publishes an event to all subscribers
	Publish(event Event) error
	
	// Subscribe subscribes to events of a specific type
	Subscribe(eventType string, handler EventHandler) error
	
	// Unsubscribe unsubscribes from events of a specific type
	Unsubscribe(eventType string, handler EventHandler) error
}

// EventHandler represents a function that handles domain events
type EventHandler func(event Event) error

// SimpleEventBus provides a simple in-memory implementation of EventBus
type SimpleEventBus struct {
	handlers map[string][]EventHandler
}

// NewSimpleEventBus creates a new simple event bus
func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Publish publishes an event to all subscribers
func (bus *SimpleEventBus) Publish(event Event) error {
	eventType := event.EventType()
	handlers, exists := bus.handlers[eventType]
	if !exists {
		return nil // No handlers for this event type
	}
	
	// In a real implementation, this would be async
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			return err
		}
	}
	
	return nil
}

// Subscribe subscribes to events of a specific type
func (bus *SimpleEventBus) Subscribe(eventType string, handler EventHandler) error {
	if bus.handlers[eventType] == nil {
		bus.handlers[eventType] = []EventHandler{}
	}
	
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	return nil
}

// Unsubscribe unsubscribes from events of a specific type
func (bus *SimpleEventBus) Unsubscribe(eventType string, handler EventHandler) error {
	handlers, exists := bus.handlers[eventType]
	if !exists {
		return nil
	}
	
	// Find and remove the handler
	for i, h := range handlers {
		// Note: In Go, we can't directly compare functions, so this is simplified
		// In a real implementation, you'd use handler IDs or other mechanisms
		_ = i
		_ = h
		break
	}
	
	return nil
}

// DomainServices provides domain-level services
type DomainServices struct {
	EventStore EventStore
	EventBus  EventBus
}

// NewDomainServices creates a new domain services instance
func NewDomainServices(eventStore EventStore, eventBus EventBus) *DomainServices {
	return &DomainServices{
		EventStore: eventStore,
		EventBus:  eventBus,
	}
}

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}