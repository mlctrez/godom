package gdutil

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestPeriodic(t *testing.T) {

	a := assert.New(t)

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	// already cancelled func should not call shouldContinue
	a.NotPanics(func() {
		Periodic(ctx, time.Millisecond, func() (ok bool) { panic("PeriodicFunc") })
	})

	ctx, cancelFunc = context.WithCancel(context.Background())
	defer cancelFunc()
	var invokeCount int32
	var shouldContinue = atomic.Bool{}
	shouldContinue.Store(true)
	go Periodic(ctx, time.Millisecond, func() (ok bool) {
		atomic.AddInt32(&invokeCount, 1)
		return shouldContinue.Load()
	})
	for atomic.LoadInt32(&invokeCount) < 2 {
		time.Sleep(10 * time.Millisecond)
	}
	lastCount := atomic.LoadInt32(&invokeCount)
	shouldContinue.Store(false)
	for atomic.LoadInt32(&invokeCount) == lastCount {
		time.Sleep(10 * time.Millisecond)
	}
	a.True(atomic.LoadInt32(&invokeCount) > 0)
}
