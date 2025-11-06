package domain

import (
	"sync"
)

// EventStore represents the interface for storing and retrieving events
type EventStore interface {
	// SaveEvent saves an event to the store
	SaveEvent(event Event) error
	
	// GetEventsByAggregate retrieves all events for a specific aggregate
	GetEventsByAggregate(aggregateID string) ([]Event, error)
	
	// GetAllEvents retrieves all events from the store
	GetAllEvents() ([]Event, error)
}

// InMemoryEventStore provides an in-memory implementation of EventStore
type InMemoryEventStore struct {
	events map[string][]Event
	mutex  sync.RWMutex
}

// NewInMemoryEventStore creates a new in-memory event store
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[string][]Event),
	}
}

// SaveEvent saves an event to the in-memory store
func (s *InMemoryEventStore) SaveEvent(event Event) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	aggregateID := event.AggregateID()
	if s.events[aggregateID] == nil {
		s.events[aggregateID] = []Event{}
	}
	
	s.events[aggregateID] = append(s.events[aggregateID], event)
	return nil
}

// GetEventsByAggregate retrieves all events for a specific aggregate
func (s *InMemoryEventStore) GetEventsByAggregate(aggregateID string) ([]Event, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	events, exists := s.events[aggregateID]
	if !exists {
		return []Event{}, nil
	}
	
	// Return a copy to prevent external modifications
	result := make([]Event, len(events))
	copy(result, events)
	return result, nil
}

// GetAllEvents retrieves all events from the store
func (s *InMemoryEventStore) GetAllEvents() ([]Event, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var allEvents []Event
	for _, events := range s.events {
		allEvents = append(allEvents, events...)
	}
	
	return allEvents, nil
}