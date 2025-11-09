package domain_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
)

// BenchmarkInMemoryEventStore_SaveEvent benchmarks event saving
func BenchmarkInMemoryEventStore_SaveEvent(b *testing.B) {
	eventStore := domain.NewInMemoryEventStore()
	event := domain.NewBaseEvent("test-id", "agg-id", "TestEvent", "test-data")

	for b.Loop() {
		_ = eventStore.SaveEvent(event)
	}
}

// BenchmarkInMemoryEventStore_GetEventsByAggregate benchmarks event retrieval
func BenchmarkInMemoryEventStore_GetEventsByAggregate(b *testing.B) {
	eventStore := domain.NewInMemoryEventStore()

	// Setup: save 1000 events for testing
	for range 1000 {
		event := domain.NewBaseEvent("test-id", "agg-id", "TestEvent", "test-data")
		_ = eventStore.SaveEvent(event)
	}

	for b.Loop() {
		_, _ = eventStore.GetEventsByAggregate("agg-id")
	}
}

// BenchmarkSimpleEventBus_Publish benchmarks event publishing with subscribers
func BenchmarkSimpleEventBus_Publish(b *testing.B) {
	eventBus := domain.NewSimpleEventBus()
	event := domain.NewBaseEvent("test-id", "agg-id", "TestEvent", "test-data")

	// Add some subscribers
	handler := func(e domain.Event) error { return nil }
	_ = eventBus.Subscribe("TestEvent", handler)
	_ = eventBus.Subscribe("TestEvent", handler)
	_ = eventBus.Subscribe("TestEvent", handler)

	for b.Loop() {
		_ = eventBus.Publish(event)
	}
}

// BenchmarkDomainServices_GenerateUUID benchmarks UUID generation
func BenchmarkDomainServices_GenerateUUID(b *testing.B) {
	for b.Loop() {
		_ = domain.GenerateUUID()
	}
}
