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
}

type TestResult struct {
	Success bool
	Message string
	Latency time.Duration
	Details map[string]interface{}
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
