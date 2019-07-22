package redis_cs

import (
	"eventsourcing/consumer"

	"github.com/go-redis/redis"
)

type RedisConsumerParams struct {
	Client           *redis.Client
	Group            string
	Consumer         string
	HandlerConsumers map[string]consumer.EventConsumerHandler
}
