package nowsm

import (
	"github.com/mlctrez/godom/app"
	"io"
	"net/http"
	"net/url"
)

var _ app.Request = (*req)(nil)

type req struct {
	r *http.Request
}

func (r *req) Method() string {
	return r.r.Method
}

func (r *req) URL() *url.URL {
	return r.r.URL
}

func (r *req) Headers() map[string]string {
	result := make(map[string]string)
	for k, v := range r.r.Header {
		result[k] = v[0]
	}
	return result
}

func (r *req) Body() io.ReadCloser {
	return r.r.Body
}

var _ app.Response = (*res)(nil)

type res struct {
	w http.ResponseWriter
}

func (r *res) SetHeader(name, value string) {
	r.w.Header().Set(name, value)
}

func (r *res) WriteHeader(statusCode int) {
	r.w.WriteHeader(statusCode)
}

func (r *res) Write(bytes []byte) (int, error) {
	return r.w.Write(bytes)
}
