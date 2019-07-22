package consumer

import (
	"context"
)

type EventConsumer interface {
	Run(ctx context.Context)
}
