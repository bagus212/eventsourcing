package redis_pbs

import (
	"eventsourcing/event"

	"github.com/go-redis/redis"
)

const (
	CHANNEL = "CHANNEL"
)

func NewRedisAdapterDomainEventPublisher(params RedisPublisherParams) RedisAdapterDomainEventPublisher {
	return RedisAdapterDomainEventPublisher{
		Channel:    params.Channel,
		Client:     params.Client,
		EventStore: params.EventStore,
	}
}

type RedisAdapterDomainEventPublisher struct {
	Channel    string
	Client     *redis.Client
	EventStore event.EventStore
}

func (adapter RedisAdapterDomainEventPublisher) Publish(event event.Event) error {
	values := map[string]interface{}{
		"event_id":   event.EventID,
		"event_name": event.EventName,
		"event_type": event.EventType,
		"data_id":    event.DataID,
		"data_name":  event.DataName,
		"data":       event.Data,
		"create_by":  event.CreateBy,
		"timestamp":  event.TimeStamp,
	}
	if err := adapter.Client.XAdd(&redis.XAddArgs{
		Stream: adapter.Channel,
		ID:     "*",
		Values: values,
	}).Err(); err != nil {
		return err
	}
	return nil
}

func (adapter RedisAdapterDomainEventPublisher) PublishAndStore(event event.Event) error {
	if err := adapter.store(event); err != nil {
		return err
	}
	return adapter.Publish(event)
}

func (adapter RedisAdapterDomainEventPublisher) store(event event.Event) error {
	return adapter.EventStore.Save(event)
}
