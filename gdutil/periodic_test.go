//go:build !wasm

package gdutil

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	shouldContinue := true
	invokeCount := 0
	go Periodic(ctx, time.Millisecond, func() (ok bool) {
		fmt.Println("invokeCount", invokeCount)
		invokeCount++
		return shouldContinue
	})
	for invokeCount < 2 {
		time.Sleep(10 * time.Millisecond)
	}
	shouldContinue = false
	time.Sleep(10 * time.Millisecond)
	a.True(invokeCount > 0)
}
