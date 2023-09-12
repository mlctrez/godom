//go:build !wasm

package godom

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConsole_logging_noWasm(t *testing.T) {
	con := Console()
	con.Log(fmt.Sprintf("Console().Log(%q)", t.Name()))
	con.Info(fmt.Sprintf("Console().Info(%q)", t.Name()))
	con.Error(fmt.Sprintf("Console().Error(%q)", t.Name()))
	con.Debug(fmt.Sprintf("Console().Debug(%q)", t.Name()))
}

func TestConsole_setup_noWasm(t *testing.T) {
	a := require.New(t)

	a.True(Global().Get("console").Truthy())

	var con *console
	var ok bool
	con, ok = Console().(*console)
	a.True(ok)
	a.NotNil(con)
	a.NotNil(con.val)

	var val *value
	val, ok = con.val.(*value)
	a.True(ok)
	a.NotNil(val)
	a.Equal(TypeObject, val.t)
	a.IsType(&consoleValue{}, val.v)

	var conVal *consoleValue
	conVal, ok = val.v.(*consoleValue)
	a.True(ok)
	a.NotNil(conVal)

	// don't change the global value which may be used by other tests
	// create a new one here to check redirectFunc delegation
	mc := &mockConsoleValue{}
	con = &console{val: mc}

	testCases := []struct {
		fn func(args ...any)
		m  string
	}{
		{fn: con.Log, m: "log"},
		{fn: con.Info, m: "info"},
		{fn: con.Error, m: "error"},
		{fn: con.Debug, m: "debug"},
	}
	for _, tc := range testCases {
		t.Run(tc.m, func(t *testing.T) {
			tc.fn("a", "b", "c")
			a.Equal(tc.m, mc.calledMethod)
			a.EqualValues([]any{"a", "b", "c"}, mc.calledArgs)
		})
	}

}

type mockConsoleValue struct {
	value
	calledMethod string
	calledArgs   []interface{}
}

func (m *mockConsoleValue) Call(method string, args ...interface{}) Value {
	m.calledMethod = method
	m.calledArgs = args
	return &value{t: TypeUndefined}
}
