//go:build js && wasm

package gfet

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	a := assert.New(t)

	headers := map[string]string{
		"custom-header": "custom-header-value",
	}
	req := &Request{URL: "/", Headers: headers}
	res, err := req.Fetch()
	a.Nil(err)
	a.True(res.Ok)

}

func TestFetch_rejected(t *testing.T) {
	a := assert.New(t)
	req := &Request{URL: "bad://url"}
	res, err := req.Fetch()
	a.Nil(res)
	a.NotNil(err)
}

var _ http.Handler = (*SlowServer)(nil)

type SlowServer struct {
}

func (s *SlowServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("got request")
	time.Sleep(5)
}

func TestFetch_signal(t *testing.T) {

	t.Run("user_abort", func(t *testing.T) {
		r := require.New(t)
		signal, cancelFunc := context.WithTimeout(context.TODO(), time.Millisecond*100)
		defer cancelFunc()
		req := &Request{URL: "/delay/150ms", Signal: signal, Headers: map[string]string{"wbt-delay-ms": "150"}}
		res, err := req.Fetch()
		r.Nil(res)
		r.NotNil(err)
		r.ErrorContains(err, "The user aborted a request.")
	})

	t.Run("no_abort", func(t *testing.T) {
		r := require.New(t)
		signal, cancelFunc := context.WithTimeout(context.TODO(), time.Millisecond*500)
		defer cancelFunc()
		req := &Request{URL: "/delay/10ms", Signal: signal, Headers: map[string]string{"wbt-delay-ms": "150"}}
		res, err := req.Fetch()
		r.NotNil(res)
		r.Nil(err)
	})

}

func TestFetch_optionsMap(t *testing.T) {
	a := assert.New(t)
	m, err := optionsMap(&Request{})
	a.Nil(err)
	a.NotNil(m)
	a.Equal(0, len(m))

	var keepAlive bool

	m, err = optionsMap(&Request{
		Method:         MethodGet,
		Headers:        map[string]string{"custom-header": "custom-value"},
		Mode:           ModeSameOrigin,
		Credentials:    CredentialsOmit,
		Cache:          CacheDefault,
		Redirect:       RedirectFollow,
		Referrer:       ReferrerNone,
		ReferrerPolicy: ReferrerPolicyNone,
		Integrity:      "Integrity",
		KeepAlive:      &keepAlive,
	})
	a.Nil(err)
	a.NotNil(m)
	a.Equal(10, len(m))

	m, err = optionsMap(&Request{Body: &BadReader{}})
	a.Nil(m)
	a.ErrorContains(err, "BadReader")

	m, err = optionsMap(&Request{Body: bytes.NewBufferString("abc")})
	a.Nil(err)
	a.Equal("abc", m["body"])

	_, err = (&Request{Body: &BadReader{}}).Fetch()
	a.ErrorContains(err, "BadReader")

}

type BadReader struct{}

func (br *BadReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("BadReader")
}
