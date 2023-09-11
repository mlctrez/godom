package gdutil

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPeriodic(t *testing.T) {

	a := assert.New(t)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	cancelFunc()
	// already cancelled func should not call shouldContinue
	a.NotPanics(func() {
		Periodic(ctx, time.Millisecond, func() (ok bool) { panic("PeriodicFunc") })
	})

	ctx, cancelFunc = context.WithCancel(context.TODO())
	defer cancelFunc()
	shouldContinue := true
	invokeCount := 0
	go Periodic(ctx, time.Millisecond, func() (ok bool) {
		invokeCount++
		return shouldContinue
	})
	time.Sleep(2 * time.Millisecond)
	shouldContinue = false
	time.Sleep(2 * time.Millisecond)
	a.True(invokeCount > 0)
}
