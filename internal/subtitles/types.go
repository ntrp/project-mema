package subtitles

import (
	"net/http"

	"media-manager/internal/subtitles/providercore"
)

type SettingValue = providercore.SettingValue
type CommandRunner = providercore.CommandRunner
type Config = providercore.Config
type MockSubtitle = providercore.MockSubtitle
type TestResult = providercore.TestResult
type SearchRequest = providercore.SearchRequest
type MediaContext = providercore.MediaContext
type MediaAlias = providercore.MediaAlias
type EpisodeNumbering = providercore.EpisodeNumbering
type FileContext = providercore.FileContext
type ReleaseProvenance = providercore.ReleaseProvenance
type Candidate = providercore.Candidate
type Download = providercore.Download

type Service struct {
	client *http.Client
}

func NewService(client *http.Client) *Service {
	if client == nil {
		client = http.DefaultClient
	}
	return &Service{client: client}
}
