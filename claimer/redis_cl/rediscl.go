package redis_cl

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
)

func NewEventFailureRecover(params RedisAdapterEventClaimer) RedisAdapterEventClaimer {
	return RedisAdapterEventClaimer{
		Client:   params.Client,
		Streams:  params.Streams,
		Group:    params.Group,
		Consumer: params.Consumer,
		MinIdle:  5 * time.Minute,
	}
}

type RedisAdapterEventClaimer struct {
	Client   *redis.Client
	Streams  []string
	Group    string
	Consumer string
	MinIdle  time.Duration
}

func (adapter RedisAdapterEventClaimer) Recover(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, stream := range adapter.Streams {
				pendingExts, err := adapter.Client.XPendingExt(&redis.XPendingExtArgs{
					Stream: stream,
					Group:  adapter.Group,
					Start:  "-",
					End:    "+",
					Count:  1000,
				}).Result()
				if err != nil {
					log.Println(err)
				}
				ids := []string{}
				for _, pendingExt := range pendingExts {
					if pendingExt.Consumer != adapter.Consumer {
						ids = append(ids, pendingExt.ID)
					}
				}
				if len(ids) > 0 {
					if err := adapter.Client.XClaim(&redis.XClaimArgs{
						Stream:   stream,
						Group:    adapter.Group,
						MinIdle:  adapter.MinIdle,
						Consumer: adapter.Consumer,
						Messages: ids,
					}).Err(); err != nil {
						log.Println(err)
					}
				}
			}
			time.Sleep(adapter.MinIdle)
		}
	}
}
