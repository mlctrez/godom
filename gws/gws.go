package gws

import (
	"fmt"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/convert"
	"net/url"
	"strings"
)

var _ WebSocket = (*webSocket)(nil)

type webSocket struct {
	v      godom.Value
	binary func(message []byte)
	text   func(message string)
}

func wsUrl(url string) string {
	replaceParts := []string{"http://", "ws://", "https://", "wss://"}
	return strings.NewReplacer(replaceParts...).Replace(url)
}

func defaultText(_ string)   {}
func defaultBinary(_ []byte) {}

func New(url string, protocols ...string) WebSocket {
	p := convert.StringsAny(protocols...)
	value := godom.Global().Value().Get("WebSocket").New(wsUrl(url), p)
	value.Set("binaryType", "arraybuffer")

	imp := &webSocket{v: value, text: defaultText, binary: defaultBinary}
	imp.addEventListener("message", imp.message)
	return imp
}

func (ws *webSocket) message(event godom.Value) {
	data := event.Get("data")
	if data.Type() == godom.TypeString {
		ws.text(data.String())
	} else {
		ws.binary(data.Bytes())
	}
}

func (ws *webSocket) OnOpen(fn godom.OnEvent) {
	ws.addEventListener("open", fn)
}

func (ws *webSocket) OnError(fn godom.OnEvent) {
	ws.addEventListener("error", fn)
}

func (ws *webSocket) OnBinaryMessage(f func(message []byte)) {
	ws.binary = f
}

func (ws *webSocket) OnTextMessage(f func(message string)) {
	ws.text = f
}

func (ws *webSocket) SendBinary(message []byte) error {
	return ws.sendBinary(message)
}

func (ws *webSocket) SendText(message string) error {
	return ws.sendText(message)
}

func (ws *webSocket) OnClose(fn func(event CloseEvent)) {
	ws.addEventListener("close", func(v godom.Value) {
		ce := CloseEvent{
			Code:     uint16(v.Get("code").Int()),
			Reason:   v.Get("reason").String(),
			WasClean: v.Get("wasClean").Bool(),
		}
		fn(ce)
	})
}

func (ws *webSocket) Close() {
	ws.v.Call("close")
}

// Rel returns a relative url using window.location.href
func Rel(p string) string {
	href := godom.Global().Location().Href()
	u, _ := url.Parse(href)
	return fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, p)
}
