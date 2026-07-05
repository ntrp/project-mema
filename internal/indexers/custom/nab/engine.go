package nab

import (
	"context"

	"media-manager/internal/indexers/engine"
)

type Engine struct {
	name   string
	client engine.HTTPDoer
}

func New(name string, client engine.HTTPDoer) *Engine {
	return &Engine{name: name, client: client}
}

func (e *Engine) Test(ctx context.Context, config engine.Config) engine.TestResult {
	return e.testCaps(ctx, config)
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	return e.search(ctx, config, query, mediaType)
}

func (e *Engine) Recent(ctx context.Context, config engine.Config) ([]engine.Release, error) {
	return e.recent(ctx, config)
}
