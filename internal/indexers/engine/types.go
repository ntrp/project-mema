package engine

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Config struct {
	ID             string
	DefinitionID   string
	Name           string
	Implementation string
	Protocol       string
	BaseURL        string
	APIKey         *string
	Categories     []int32
	Fields         json.RawMessage
	Redirect       bool
}

type TestResult struct {
	Success bool
	Message string
	Latency time.Duration
	Details map[string]interface{}
}

type Release struct {
	IndexerID       string
	IndexerName     string
	IndexerProtocol string
	Title           string
	DownloadURL     string
	InfoURL         string
	GUID            string
	SizeBytes       int64
	Seeders         *int32
	Peers           *int32
	PublishedAt     *time.Time
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Engine interface {
	Test(ctx context.Context, config Config) TestResult
	Search(ctx context.Context, config Config, query string, mediaType string) ([]Release, error)
}
