package redis_pbs

import (
	"eventsourcing/event"

	"github.com/go-redis/redis"
)

type RedisPublisherParams struct {
	Channel    string
	Client     *redis.Client
	EventStore event.EventStore
}
