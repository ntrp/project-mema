package gazelleapi

import (
	"context"
	"net/http"

	"media-manager/internal/indexers/engine"
)

type Engine struct {
	client  engine.HTTPDoer
	options Options
}

func New(options Options, clients ...engine.HTTPDoer) *Engine {
	var client engine.HTTPDoer
	if len(clients) > 0 {
		client = clients[0]
	}
	return &Engine{client: client, options: options}
}

func (e *Engine) Test(ctx context.Context, config engine.Config) engine.TestResult {
	if e.client == nil {
		return engine.FailedResult("Connection failed", "error", "HTTP client is not configured")
	}
	endpoint, err := searchURL(config, e.options, "test")
	if err != nil {
		return engine.FailedResult("Invalid "+e.options.Name+" request", "error", err.Error())
	}
	req, err := authorizedRequest(ctx, e.client, config, e.options, endpoint)
	if err != nil {
		return engine.FailedResult("Invalid "+e.options.Name+" request", "error", err.Error())
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return engine.RequestFailedResult(err)
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.StatusFailedResult(resp.StatusCode)
	}
	return engine.SuccessResult(e.options.Name+" indexer reachable", "endpoint", endpoint)
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, e.options, query)
	if err != nil {
		return nil, err
	}
	req, err := authorizedRequest(ctx, e.client, config, e.options, endpoint)
	if err != nil {
		return nil, err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, engine.HTTPStatusError(resp)
	}
	body, err := engine.ReadLimitedBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseReleases(config, e.options, body)
}
