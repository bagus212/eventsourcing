package recover

import (
	"context"
)

type EventClaimer interface {
	Claim(ctx context.Context)
}
