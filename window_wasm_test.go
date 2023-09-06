package godom

import (
	"github.com/stretchr/testify/assert"
	"syscall/js"
	"testing"
)

func TestLocation_Reload(t *testing.T) {
	a := assert.New(t)
	locValue := js.Global().Get("Object").New()
	reloadCalled := false
	locValue.Set("reload", js.FuncOf(func(this js.Value, args []js.Value) any {
		reloadCalled = true
		return nil
	}))
	l := &location{this: FromJsValue(locValue)}
	l.Reload()
	a.True(reloadCalled)
}
