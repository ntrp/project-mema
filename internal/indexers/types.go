package indexers

import (
	"net/http"
	"time"
)

type Config struct {
	ID         string
	Name       string
	Type       string
	BaseURL    string
	APIKey     *string
	Categories []int32
}

type TestResult struct {
	Success bool
	Message string
	Latency time.Duration
	Details map[string]interface{}
}

type Release struct {
	IndexerID   string
	IndexerName string
	IndexerType string
	Title       string
	DownloadURL string
	InfoURL     string
	GUID        string
	SizeBytes   int64
	Seeders     *int32
	Peers       *int32
	PublishedAt *time.Time
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
