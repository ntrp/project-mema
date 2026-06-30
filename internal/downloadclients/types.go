package downloadclients

import (
	"net/http"
	"time"
)

type Config struct {
	Name     string
	Type     string
	BaseURL  string
	Username *string
	Password *string
	APIKey   *string
	Category *string
}

type TestResult struct {
	Success bool
	Message string
	Latency time.Duration
	Details map[string]interface{}
}

type AddRequest struct {
	URL      string
	Title    string
	Category *string
}

type AddResult struct {
	Success    bool
	Message    string
	DownloadID string
	Details    map[string]interface{}
}

type StatusRequest struct {
	DownloadID string
}

type CancelRequest struct {
	DownloadID string
}

type CancelResult struct {
	Success bool
	Message string
	Details map[string]interface{}
}

type StatusFile struct {
	Path      string
	SizeBytes int64
	Complete  bool
}

type StatusResult struct {
	Success         bool
	Found           bool
	Status          string
	ProgressPercent *int
	Message         string
	Files           []StatusFile
	Details         map[string]interface{}
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
