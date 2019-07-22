package publisher

import "eventsourcing/event"

type DomainEventPublisher interface {
	Publish(event event.Event) error
	PublishAndStore(event event.Event) error
}
