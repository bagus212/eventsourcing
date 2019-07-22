package redis_cs

import (
	"context"
	"encoding/json"
	"eventsourcing/consumer"
	"eventsourcing/event"
	"log"

	"github.com/go-redis/redis"
)

func NewRedisAdapterEventConsumer(params RedisConsumerParams) (RedisAdaptersEventConsumer, error) {
	var streams []string
	for key := range params.HandlerConsumers {
		if err := params.Client.XGroupCreateMkStream(key, params.Group, "0").Err(); err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
			return RedisAdaptersEventConsumer{}, err
		}
		streams = append(streams, key)
	}
	for range params.HandlerConsumers {
		streams = append(streams, ">")
	}
	return RedisAdaptersEventConsumer{
		Client:           params.Client,
		HandlerConsumers: params.HandlerConsumers,
		Streams:          streams,
		Group:            params.Group,
		Consumer:         params.Consumer,
	}, nil
}

type RedisAdaptersEventConsumer struct {
	Client           *redis.Client
	HandlerConsumers map[string]consumer.EventConsumerHandler
	Streams          []string
	Group            string
	Consumer         string
}

func (adapter *RedisAdaptersEventConsumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			streams, err := adapter.Client.XReadGroup(&redis.XReadGroupArgs{
				Streams:  adapter.Streams,
				Group:    adapter.Group,
				Consumer: adapter.Consumer,
				NoAck:    false,
				Count:    100,
				Block:    0,
			}).Result()
			if err != nil {
				log.Println(err)
				continue
			}
			for _, stream := range streams {
				for _, message := range stream.Messages {
					event, err := adapter.messageToEvent(message.Values)
					if err != nil {
						log.Println(err)
						continue
					}
					err = adapter.HandlerConsumers[stream.Stream].Apply(event)
					if err != nil {
						log.Println(err)
						continue
					}
					err = adapter.Client.XAck(stream.Stream, adapter.Group, message.ID).Err()
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}

		}

	}
}

func (adapter *RedisAdaptersEventConsumer) messageToEvent(values map[string]interface{}) (event.Event, error) {
	var event event.Event
	valuesString, err := json.Marshal(values)
	if err != nil {
		return event, err
	}
	err = json.Unmarshal([]byte(valuesString), &event)
	if err != nil {
		return event, err
	}
	return event, nil
}
