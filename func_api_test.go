package godom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockFunc struct {
	this   Value
	args   []Value
	result any
}

func (m *mockFunc) target(this Value, args []Value) any {
	m.this = this
	m.args = args
	return m.result
}

func TestFunc_Call(t *testing.T) {
	a := assert.New(t)
	a.NotNil(a)

	mock := &mockFunc{}
	fn := FuncOf(mock.target)
	defer fn.Release()

	g := Global()
	g.Set("valueForTestFuncOf", fn)
	g.Call("valueForTestFuncOf", "a", "b")

	a.NotNil(mock.this)
	a.NotNil(mock.args)
	a.Equal(2, len(mock.args))
	a.Equal("a", mock.args[0].String())
}

func TestFunc_Invoke(t *testing.T) {
	a := assert.New(t)
	a.NotNil(a)

	mock := &mockFunc{}
	fn := FuncOf(mock.target)
	defer fn.Release()

	g := Global()
	g.Set("valueForTestFuncOf", fn)

	g.Get("valueForTestFuncOf").Invoke("a", "b")

	a.NotNil(mock.this)
	a.NotNil(mock.args)
	a.Equal(2, len(mock.args))
	a.Equal("a", mock.args[0].String())
}
