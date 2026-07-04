package unsupported

import (
	"context"
	"fmt"

	"media-manager/internal/indexers/engine"
)

type Engine struct {
	name string
}

func New(name string) *Engine {
	return &Engine{name: name}
}

func (e *Engine) Test(ctx context.Context, config engine.Config) engine.TestResult {
	return engine.FailedResult(
		"Custom indexer engine is not implemented",
		"implementation", e.name,
		"definitionId", config.DefinitionID,
	)
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	return nil, fmt.Errorf("custom indexer engine %q is not implemented", e.name)
}
