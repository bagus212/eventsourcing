package redis_cl

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisClaimerParams struct {
	Client   *redis.Client
	Streams  []string
	Group    string
	Consumer string
	MinIdle  time.Duration
}
