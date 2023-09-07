package godom

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestFuncOf(t *testing.T) {
	a := assert.New(t)
	a.NotNil(a)

	var thisInvoked Value
	var argsInvoked []Value
	fn := FuncOf(func(this Value, args []Value) any {
		thisInvoked = this
		argsInvoked = args
		return nil
	})
	defer fn.Release()

	// for now this test is only good on js
	if runtime.GOOS == "js" {
		global := Global().Value()
		global.Set("valueForTestFuncOf", fn)
		global.Call("valueForTestFuncOf", "a", "b")
		a.NotNil(thisInvoked)
		a.NotNil(argsInvoked)
		a.Equal("a", argsInvoked[0].String())
	}

}
