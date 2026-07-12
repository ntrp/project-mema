package subtitles

import (
	"context"
	"net/http"
	"time"
)

type SettingValue struct {
	StringValue  *string
	NumberValue  *float64
	BooleanValue *bool
	StringValues []string
}

type CommandRunner func(ctx context.Context, name string, args ...string) ([]byte, error)

type Config struct {
	Name           string
	Type           string
	BaseURL        string
	Username       *string
	Password       *string
	APIKey         *string
	Settings       map[string]SettingValue
	SecretSettings map[string]string
	CommandRunner  CommandRunner
	MockSubtitles  []MockSubtitle
}

type MockSubtitle struct {
	Title      string
	LanguageID string
	Format     string
}

type TestResult struct {
	Success bool
	Message string
	Latency time.Duration
	Details map[string]any
}

type SearchRequest struct {
	MediaType     string
	Title         string
	LanguageID    string
	Year          *int32
	SeasonNumber  *int32
	EpisodeNumber *int32
	FilePath      string
}

type Candidate struct {
	ProviderName  string
	LanguageID    string
	FileID        int64
	Format        string
	ReleaseName   string
	DownloadCount int
	SourceURL     string
	SourceRef     string
}

type Download struct {
	Content []byte
	URL     string
}

type Service struct {
	client *http.Client
}

func NewService(client *http.Client) *Service {
	if client == nil {
		client = http.DefaultClient
	}
	return &Service{client: client}
}
