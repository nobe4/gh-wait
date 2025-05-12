package looper

import (
	"context"
)

type Looper interface {
	Loop(ctx context.Context) error
	Stop()
}
