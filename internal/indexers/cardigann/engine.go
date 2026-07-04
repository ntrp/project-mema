package cardigann

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"media-manager/internal/indexers/engine"
)

type Config = engine.Config
type HTTPDoer = engine.HTTPDoer
type Release = engine.Release
type TestResult = engine.TestResult

type Engine struct {
	client engine.HTTPDoer
	loader *cardigannLoader
}

func New(client engine.HTTPDoer) *Engine {
	return &Engine{client: client, loader: newCardigannLoader(client)}
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	return e.searchCardigann(ctx, config, query, mediaType)
}

func (e *Engine) Test(ctx context.Context, config engine.Config) engine.TestResult {
	return e.testCardigann(ctx, config)
}

func (e *Engine) UseLocalDefinitions(definitions map[string]string) {
	e.loader.remote = ""
	e.loader.local = definitions
}

func closeBody(body io.ReadCloser) {
	engine.CloseBody(body)
}

func readLimitedBody(body io.Reader) ([]byte, error) {
	return engine.ReadLimitedBody(body)
}

func parseBaseURL(baseURL string) (*url.URL, error) {
	return engine.ParseBaseURL(baseURL)
}

func stringValue(value *string) string {
	return engine.StringValue(value)
}

func firstNonEmpty(values ...string) string {
	return engine.FirstNonEmpty(values...)
}

func httpStatusError(resp *http.Response) error {
	return engine.HTTPStatusError(resp)
}

func failedResult(message string, pairs ...interface{}) TestResult {
	return engine.FailedResult(message, pairs...)
}

func successResult(message string, pairs ...interface{}) TestResult {
	return engine.SuccessResult(message, pairs...)
}

func requestFailedResult(err error) TestResult {
	return engine.RequestFailedResult(err)
}

func statusFailedResult(statusCode int) TestResult {
	return engine.StatusFailedResult(statusCode)
}
