package gdutil

import (
	"context"
	"time"
)

type PeriodicFunc func() (ok bool)

// Periodic executes a PeriodicFunc every interval.
// It will return if ctx is cancelled or PeriodicFunc returns false.
// It is intended to be run as a go routine.
//
//	go Periodic(ctx, time.Second, func() bool {
//	  err := ErrorsWhenDone()
//	  // no error means continue
//	  return err == nil
//	})
func Periodic(ctx context.Context, interval time.Duration, pf PeriodicFunc) {
	pingTicker := time.NewTicker(interval)
	defer pingTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			if pf() {
				continue
			}
			return
		}
	}
}
