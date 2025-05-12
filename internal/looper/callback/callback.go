package callback

import (
	"context"
	"time"
)

type Looper struct {
	cb     func()
	delay  time.Duration
	stopCh chan struct{}
}

func New(cb func(), delay time.Duration) Looper {
	return Looper{
		cb:     cb,
		delay:  delay,
		stopCh: make(chan struct{}, 1),
	}
}

func (c Looper) Loop(ctx context.Context) error {
	ticker := time.NewTicker(c.delay)
	defer ticker.Stop()

	// TODO add that back later
	// timer := time.NewTimer(10 * time.Second)
	// defer timer.Stop()

	c.cb()

	for {
		select {
		case <-ticker.C:
			c.cb()

		case <-ctx.Done():
			return ctx.Err()
		// case <-timer.C:
		// 	return errors.New("timer expired")

		case <-c.stopCh:
			return nil
		}
	}
}

func (c Looper) Stop() {
	c.stopCh <- struct{}{}
}
