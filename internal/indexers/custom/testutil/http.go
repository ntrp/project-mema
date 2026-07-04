package testutil

import (
	"io"
	"net/http"
	"strings"
)

type FakeHTTPDoer func(req *http.Request) (*http.Response, error)

func (f FakeHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func Response(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
