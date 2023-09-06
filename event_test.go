//go:build !js && !wasm

package godom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventFunc(t *testing.T) {
	a := assert.New(t)
	gotEvent := false
	eventFunc := EventFunc(func() { gotEvent = true })
	eventFunc(ToValue("test"))
	a.True(gotEvent)
}
