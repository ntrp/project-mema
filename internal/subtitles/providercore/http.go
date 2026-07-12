package providercore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"media-manager/internal/subtitles/security"
)

type HTTPClient struct {
	providerKey string
	client      *http.Client
	headers     http.Header
	download    bool
}

func NewHTTPClient(providerKey string, base *http.Client) *HTTPClient {
	return newHTTPClient(providerKey, base, false)
}

func NewDownloadHTTPClient(providerKey string, base *http.Client) *HTTPClient {
	return newHTTPClient(providerKey, base, true)
}

func newHTTPClient(providerKey string, base *http.Client, download bool) *HTTPClient {
	if base == nil {
		base = http.DefaultClient
	}
	clone := *base
	previousRedirect := clone.CheckRedirect
	clone.CheckRedirect = func(request *http.Request, via []*http.Request) error {
		if len(via) > 0 {
			fromURL := via[len(via)-1].URL.String()
			if err := security.ValidateRedirect(providerKey, fromURL, request.URL.String(), download); err != nil {
				return err
			}
		}
		if previousRedirect != nil {
			return previousRedirect(request, via)
		}
		return nil
	}
	if clone.Jar == nil {
		clone.Jar, _ = cookiejar.New(nil)
	}
	if clone.Timeout == 0 {
		clone.Timeout = 30 * time.Second
	}
	return &HTTPClient{providerKey: providerKey, client: &clone, headers: http.Header{}, download: download}
}

func (c *HTTPClient) SetHeader(key string, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	c.headers.Set(key, value)
}

func (c *HTTPClient) Request(ctx context.Context, method string, rawURL string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	for key, values := range c.headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	return request, nil
}

func (c *HTTPClient) Do(request *http.Request, maxBytes int64) ([]byte, *http.Response, error) {
	if request == nil || request.URL == nil {
		return nil, nil, fmt.Errorf("request URL is required")
	}
	if err := security.ValidateProviderURL(c.providerKey, request.URL.String(), c.download); err != nil {
		return nil, nil, err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return nil, response, err
	}
	defer response.Body.Close()
	return readHTTPResponse(response, maxBytes)
}

func (c *HTTPClient) Get(ctx context.Context, rawURL string, maxBytes int64) ([]byte, *http.Response, error) {
	request, err := c.Request(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, nil, err
	}
	return c.Do(request, maxBytes)
}

func readHTTPResponse(response *http.Response, maxBytes int64) ([]byte, *http.Response, error) {
	if maxBytes <= 0 {
		maxBytes = 10 << 20
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, maxBytes+1))
	if err != nil {
		return nil, response, err
	}
	if int64(len(data)) > maxBytes {
		return nil, response, fmt.Errorf("response size limit exceeded")
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, response, fmt.Errorf("http status %d", response.StatusCode)
	}
	return data, response, nil
}
