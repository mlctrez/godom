package godom

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"syscall/js"
	"testing"
)

func TestConsole_WithWasm(t *testing.T) {
	a := require.New(t)

	con, ok := Console().(*console)
	a.True(ok)
	a.IsType(&console{}, con)
	a.NotNil(con.val)
	jsConsole := js.Global().Get("console")
	a.True(jsConsole.Equal(ToJsValue(con.val).(js.Value)))

	// fire off some test logging
	con.Log(fmt.Sprintf("Console().Log(%q)", t.Name()))
	con.Info(fmt.Sprintf("Console().Info(%q)", t.Name()))
	con.Error(fmt.Sprintf("Console().Error(%q)", t.Name()))
	con.Debug(fmt.Sprintf("Console().Debug(%q)", t.Name()))
}
