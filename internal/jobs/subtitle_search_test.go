package jobs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

func TestSubtitleSearchRequestBuildsEpisodeContext(t *testing.T) {
	year := int32(2026)
	item := storage.MediaItem{
		ID:                    uuid.New(),
		Type:                  "serie",
		Title:                 "Scenario Series",
		Year:                  &year,
		FilePaths:             []string{"/library/Scenario.Series.S01E02.mkv"},
		SubtitlePreferredMode: "mixed",
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{
			{LanguageID: "english"},
		},
	}

	request, ok := subtitleSearchRequest(item, SubtitleSearchArgs{})

	if !ok {
		t.Fatal("expected subtitle search request")
	}
	if request.Title != "Scenario Series" || request.LanguageID != "english" {
		t.Fatalf("request = %#v", request)
	}
	if request.SeasonNumber == nil || *request.SeasonNumber != 1 {
		t.Fatalf("season = %#v", request.SeasonNumber)
	}
	if request.EpisodeNumber == nil || *request.EpisodeNumber != 2 {
		t.Fatalf("episode = %#v", request.EpisodeNumber)
	}
}

func TestSubtitleSearchDownloadsAndRecordsSubtitle(t *testing.T) {
	ctx, store := jobsTestStore(t)
	tmp := t.TempDir()
	mediaPath := filepath.Join(tmp, "Scenario.Movie.2026.mkv")
	if err := os.WriteFile(mediaPath, []byte("movie"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := subtitleProviderServer(t, http.StatusOK)
	apiKey := "subtitle-key"
	provider, err := store.CreateSubtitleProvider(ctx, storage.SubtitleProviderInput{
		Name:    "Scenario Subtitles",
		Type:    "opensubtitles",
		BaseURL: server.URL,
		APIKey:  &apiKey,
		Enabled: true,
	})
	if err != nil || provider.ID == uuid.Nil {
		t.Fatalf("provider = %#v err=%v", provider, err)
	}
	item, err := store.CreateMediaItem(ctx, storage.MediaItemInput{
		Type:      "movie",
		Title:     "Scenario Movie",
		Year:      int32Ptr(2026),
		Monitored: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	item.FilePaths = []string{mediaPath}
	item.SubtitlePreferredMode = "mixed"
	item.SubtitleTargets = []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english"},
	}

	err = subtitleSearchDownload(ctx, store, subtitles.NewService(server.Client()), nil, item, SubtitleSearchArgs{LanguageID: "english"})

	if err != nil {
		t.Fatalf("subtitleSearchDownload returned error: %v", err)
	}
	subtitlePath := filepath.Join(tmp, "Scenario.Movie.2026.english.srt")
	content, err := os.ReadFile(subtitlePath)
	if err != nil {
		t.Fatalf("read subtitle: %v", err)
	}
	if string(content) != "1\n00:00:00,000 --> 00:00:01,000\nScenario\n" {
		t.Fatalf("subtitle content = %q", content)
	}
	records, err := store.ListMediaItemSubtitles(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].LanguageID != "english" || records[0].FilePath != subtitlePath {
		t.Fatalf("records = %#v", records)
	}
	events, err := store.ListSystemEvents(ctx, 20, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !hasSystemEvent(events, "Subtitle downloaded") {
		t.Fatalf("events = %#v", events)
	}
}

func TestSubtitleSearchConvertsDownloadedSubtitleToTargetFormat(t *testing.T) {
	ctx, store := jobsTestStore(t)
	tmp := t.TempDir()
	mediaPath := filepath.Join(tmp, "Scenario.Movie.2026.mkv")
	if err := os.WriteFile(mediaPath, []byte("movie"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := subtitleProviderServer(t, http.StatusOK)
	apiKey := "subtitle-key"
	if _, err := store.CreateSubtitleProvider(ctx, storage.SubtitleProviderInput{
		Name: "Scenario Subtitles", Type: "opensubtitles", BaseURL: server.URL, APIKey: &apiKey, Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, storage.MediaItemInput{
		Type: "movie", Title: "Scenario Movie", Year: int32Ptr(2026), Monitored: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	item.FilePaths = []string{mediaPath}
	item.SubtitlePreferredMode = "mixed"
	item.SubtitleTargets = []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english", Formats: []string{"vtt"}},
	}

	err = subtitleSearchDownload(ctx, store, subtitles.NewService(server.Client()), nil, item, SubtitleSearchArgs{LanguageID: "english"})

	if err != nil {
		t.Fatalf("subtitleSearchDownload returned error: %v", err)
	}
	subtitlePath := filepath.Join(tmp, "Scenario.Movie.2026.english.vtt")
	content, err := os.ReadFile(subtitlePath)
	if err != nil {
		t.Fatalf("read subtitle: %v", err)
	}
	if !strings.HasPrefix(string(content), "WEBVTT") {
		t.Fatalf("subtitle content = %q", content)
	}
	records, err := store.ListMediaItemSubtitles(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Format != "vtt" || records[0].FilePath != subtitlePath {
		t.Fatalf("records = %#v", records)
	}
}

func TestSubtitleSearchRejectsUnsupportedBitmapTargetFormat(t *testing.T) {
	ctx, store := jobsTestStore(t)
	tmp := t.TempDir()
	mediaPath := filepath.Join(tmp, "Scenario.Movie.2026.mkv")
	if err := os.WriteFile(mediaPath, []byte("movie"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := subtitleProviderServer(t, http.StatusOK)
	apiKey := "subtitle-key"
	if _, err := store.CreateSubtitleProvider(ctx, storage.SubtitleProviderInput{
		Name: "Scenario Subtitles", Type: "opensubtitles", BaseURL: server.URL, APIKey: &apiKey, Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, storage.MediaItemInput{
		Type: "movie", Title: "Scenario Movie", Year: int32Ptr(2026), Monitored: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	item.FilePaths = []string{mediaPath}
	item.SubtitlePreferredMode = "mixed"
	item.SubtitleTargets = []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english", Formats: []string{"pgs"}},
	}

	err = subtitleSearchDownload(ctx, store, subtitles.NewService(server.Client()), nil, item, SubtitleSearchArgs{LanguageID: "english"})

	if err == nil || !strings.Contains(err.Error(), "pgs") {
		t.Fatalf("expected pgs conversion error, got %v", err)
	}
}

func TestSubtitleSearchProviderFailureSurfacesEvent(t *testing.T) {
	ctx, store := jobsTestStore(t)
	server := subtitleProviderServer(t, http.StatusBadGateway)
	apiKey := "subtitle-key"
	if _, err := store.CreateSubtitleProvider(ctx, storage.SubtitleProviderInput{
		Name:    "Broken Subtitles",
		Type:    "opensubtitles",
		BaseURL: server.URL,
		APIKey:  &apiKey,
		Enabled: true,
	}); err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, storage.MediaItemInput{
		Type:      "movie",
		Title:     "Scenario Movie",
		Monitored: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	item.FilePaths = []string{filepath.Join(t.TempDir(), "Scenario.Movie.mkv")}

	err = subtitleSearchDownload(ctx, store, subtitles.NewService(server.Client()), nil, item, SubtitleSearchArgs{LanguageID: "english"})

	if err == nil {
		t.Fatal("expected provider failure")
	}
	events, err := store.ListSystemEvents(ctx, 20, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !hasSystemEvent(events, "Subtitle search failed") {
		t.Fatalf("events = %#v", events)
	}
}

func subtitleProviderServer(t *testing.T, searchStatus int) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Api-Key") != "subtitle-key" {
			t.Errorf("missing api key header")
		}
		switch r.URL.Path {
		case "/api/v1/subtitles":
			w.WriteHeader(searchStatus)
			if searchStatus == http.StatusOK {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"data": []map[string]any{{
						"attributes": map[string]any{
							"language":       "en",
							"download_count": 10,
							"url":            "https://provider.test/subtitle/44",
							"files": []map[string]any{{
								"file_id":   44,
								"file_name": "Scenario.Movie.2026.srt",
							}},
						},
					}},
				})
			}
		case "/api/v1/download":
			_ = json.NewEncoder(w).Encode(map[string]any{"link": "http://" + r.Host + "/file.srt"})
		case "/file.srt":
			_, _ = w.Write([]byte("1\n00:00:00,000 --> 00:00:01,000\nScenario\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func hasSystemEvent(events []storage.SystemEvent, message string) bool {
	for _, event := range events {
		if event.Message == message {
			return true
		}
	}
	return false
}
