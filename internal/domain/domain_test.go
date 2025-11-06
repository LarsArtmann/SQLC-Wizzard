package domain_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/domain"
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

var _ = Describe("Domain Events", func() {
	var (
		eventStore domain.EventStore
		eventBus  domain.EventBus
		services  *domain.DomainServices
	)

	BeforeEach(func() {
		eventStore = domain.NewInMemoryEventStore()
		eventBus = domain.NewSimpleEventBus()
		services = domain.NewDomainServices(eventStore, eventBus)
	})

	AfterEach(func() {
		// Cleanup if needed
	})

	Describe("BaseEvent", func() {
		It("should create a valid base event", func() {
			eventData := map[string]string{"test": "data"}
			event := domain.NewBaseEvent("event-123", "aggregate-456", "TestEvent", eventData)

			Expect(event.ID()).To(Equal("event-123"))
			Expect(event.AggregateID()).To(Equal("aggregate-456"))
			Expect(event.EventType()).To(Equal("TestEvent"))
			Expect(event.Data()).To(Equal(eventData))
			Expect(event.OccurredAt()).To(BeTemporally("~", time.Now(), time.Second))
		})
	})

	Describe("InMemoryEventStore", func() {
		It("should save and retrieve events by aggregate ID", func() {
			event1 := domain.NewBaseEvent("e1", "agg1", "Event1", "data1")
			event2 := domain.NewBaseEvent("e2", "agg1", "Event2", "data2")
			event3 := domain.NewBaseEvent("e3", "agg2", "Event3", "data3")

			err := eventStore.SaveEvent(event1)
			Expect(err).ToNot(HaveOccurred())

			err = eventStore.SaveEvent(event2)
			Expect(err).ToNot(HaveOccurred())

			err = eventStore.SaveEvent(event3)
			Expect(err).ToNot(HaveOccurred())

			// Retrieve events for agg1
			agg1Events, err := eventStore.GetEventsByAggregate("agg1")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(agg1Events)).To(Equal(2))
			Expect(agg1Events[0].ID()).To(Equal("e1"))
			Expect(agg1Events[1].ID()).To(Equal("e2"))

			// Retrieve events for agg2
			agg2Events, err := eventStore.GetEventsByAggregate("agg2")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(agg2Events)).To(Equal(1))
			Expect(agg2Events[0].ID()).To(Equal("e3"))
		})

		It("should return empty slice for non-existent aggregate", func() {
			events, err := eventStore.GetEventsByAggregate("non-existent")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(events)).To(Equal(0))
		})
	})

	Describe("SimpleEventBus", func() {
		It("should publish events to subscribers", func() {
			receivedEvent := false
			handler := func(event domain.Event) error {
				receivedEvent = true
				return nil
			}

			err := eventBus.Subscribe("TestEvent", handler)
			Expect(err).ToNot(HaveOccurred())

			event := domain.NewBaseEvent("e1", "agg1", "TestEvent", "data")
			err = eventBus.Publish(event)
			Expect(err).ToNot(HaveOccurred())
			Expect(receivedEvent).To(BeTrue())
		})

		It("should not publish events when no subscribers", func() {
			event := domain.NewBaseEvent("e1", "agg1", "TestEvent", "data")
			err := eventBus.Publish(event)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("DomainServices", func() {
		It("should create domain services", func() {
			Expect(services.EventStore).ToNot(BeNil())
			Expect(services.EventBus).ToNot(BeNil())
		})

		It("should generate UUIDs", func() {
			uuid1 := domain.GenerateUUID()
			uuid2 := domain.GenerateUUID()

			Expect(uuid1).ToNot(BeEmpty())
			Expect(uuid2).ToNot(BeEmpty())
			Expect(uuid1).ToNot(Equal(uuid2))
		})
	})

	Describe("Project Events", func() {
		It("should create ProjectCreated event data", func() {
			eventData := domain.ProjectCreated{
				ProjectID:   "proj-123",
				Name:        "Test Project",
				ProjectType: generated.ProjectType("microservice"),
				Database:    generated.DatabaseType("postgresql"),
				CreatedAt:   "2023-01-01T00:00:00Z",
			}

			Expect(eventData.ProjectID).To(Equal("proj-123"))
			Expect(eventData.Name).To(Equal("Test Project"))
			Expect(eventData.ProjectType).To(Equal(generated.ProjectType("microservice")))
			Expect(eventData.Database).To(Equal(generated.DatabaseType("postgresql")))
		})

		It("should create ConfigValidated event data", func() {
			eventData := domain.ConfigValidated{
				ProjectID:        "proj-123",
				IsValid:          true,
				ValidationErrors:  []string{},
				ValidatedAt:      "2023-01-01T00:00:00Z",
			}

			Expect(eventData.ProjectID).To(Equal("proj-123"))
			Expect(eventData.IsValid).To(BeTrue())
			Expect(len(eventData.ValidationErrors)).To(Equal(0))
		})
	})
})

func TestDomainEvents(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Domain Events Suite")
}