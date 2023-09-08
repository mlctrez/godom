package gws

import (
	"github.com/mlctrez/godom"
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func Test_wsUrl(t *testing.T) {
	a := assert.New(t)
	a.Equal("ws://somehost", wsUrl("http://somehost"))
	a.Equal("wss://somehost", wsUrl("https://somehost"))
}

func TestNew(t *testing.T) {
	if runtime.GOOS != "js" {
		return
	}
	a := assert.New(t)
	ws := New("http://127.0.0.1:39999")
	defer ws.Close()

	errChan := make(chan godom.Value, 1)
	ws.OnError(func(event godom.Value) { errChan <- event })

	closeChan := make(chan CloseEvent, 1)
	ws.OnClose(func(event CloseEvent) { closeChan <- event })

	ws.OnOpen(func(event godom.Value) { a.Fail("open should not be called") })

	var errEvent godom.Value
	var closeEvent CloseEvent
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
	case errEvent = <-errChan:
	}

	select {
	case <-timer.C:
	case closeEvent = <-closeChan:
	}

	a.NotNil(errEvent)
	a.Equal("error", errEvent.Get("type").String())

	a.NotNil(closeEvent)
	a.Equal(1006, int(closeEvent.Code))
}

func TestWebSocket_messageHandlers(t *testing.T) {
	a := assert.New(t)

	ws := &webSocket{}
	ws.OnBinaryMessage(defaultBinary)
	ws.OnTextMessage(defaultText)

	ptrEq := func(one, two any) bool {
		return reflect.ValueOf(one).Pointer() == reflect.ValueOf(two).Pointer()
	}

	defaultBinary(nil)
	defaultText("")

	a.True(ptrEq(defaultText, ws.text))
	a.True(ptrEq(defaultBinary, ws.binary))

	a.IsType(&webSocket{}, ws)
}

func TestRel(t *testing.T) {
	if runtime.GOOS != "js" {
		return
	}
	a := assert.New(t)
	rel := Rel("ws")
	a.True(strings.HasPrefix(rel, "http://"))
	a.True(strings.HasSuffix(rel, "/ws"))
}

func TestCloseFunc(t *testing.T) {
	a := assert.New(t)
	called := false
	closeFunc := CloseFunc(func() { called = true })
	closeFunc(CloseEvent{})
	a.True(called)
}
