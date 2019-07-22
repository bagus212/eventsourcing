package consumer

import "eventsourcing/event"

type EventConsumerHandler interface {
	Apply(event event.Event) error
}
