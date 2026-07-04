package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"media-manager/internal/storage"
)

func TestSCNSettings009IntegrationInputsNormalizeAndValidate(t *testing.T) {
	apiKey := " key "
	username := " user "
	password := " pass "
	category := " movies "
	download, ok := downloadClientInput(httptest.NewRecorder(), DownloadClientRequest{
		Name:     " Client ",
		Type:     Transmission,
		Protocol: IndexerProtocolTorrent,
		BaseUrl:  " http://client.local ",
		Username: &username,
		Password: &password,
		ApiKey:   &apiKey,
		Category: &category,
		Enabled:  true,
		Priority: 5,
	})
	if !ok || download.Name != "Client" || download.Protocol != "torrent" || download.BaseURL != "http://client.local" {
		t.Fatalf("download input = %#v, ok = %v", download, ok)
	}
	if download.APIKey == nil || *download.APIKey != "key" || download.Category == nil || *download.Category != "movies" {
		t.Fatalf("download optional fields = %#v", download)
	}

	categories := []int32{2000, 2040}
	indexer, ok := indexerInput(httptest.NewRecorder(), IndexerRequest{
		DefinitionId: "generic-torznab",
		Name:         " Indexer ",
		BaseUrl:      " http://indexer.local ",
		ApiKey:       &apiKey,
		Categories:   &categories,
		Enabled:      true,
		Priority:     10,
	}, testLanguageCatalog())
	if !ok || indexer.Name != "Indexer" || len(indexer.Categories) != 2 {
		t.Fatalf("indexer input = %#v, ok = %v", indexer, ok)
	}
	if indexer.IndexerURLs == nil || indexer.LegacyURLs == nil {
		t.Fatalf("indexer URL arrays must be normalized to empty slices: %#v", indexer)
	}

	token := " access "
	provider, ok := metadataProviderInput(httptest.NewRecorder(), MetadataProviderRequest{
		Name:        " Provider ",
		Type:        Tmdb,
		BaseUrl:     " http://provider.local ",
		ApiKey:      &apiKey,
		AccessToken: &token,
		Enabled:     true,
		Priority:    20,
	})
	if !ok || provider.Name != "Provider" || provider.AccessToken == nil || *provider.AccessToken != "access" {
		t.Fatalf("provider input = %#v, ok = %v", provider, ok)
	}

	assertBadRequest(t, func(w http.ResponseWriter) bool {
		_, ok := downloadClientInput(w, DownloadClientRequest{Type: Transmission, BaseUrl: "http://x", Priority: 1})
		return ok
	})
	assertBadRequest(t, func(w http.ResponseWriter) bool {
		_, ok := downloadClientInput(w, DownloadClientRequest{
			Name:     "Bad",
			Type:     Sabnzbd,
			Protocol: IndexerProtocolTorrent,
			BaseUrl:  "http://x",
			Priority: 1,
		})
		return ok
	})
	assertBadRequest(t, func(w http.ResponseWriter) bool {
		_, ok := indexerInput(
			w,
			IndexerRequest{DefinitionId: "bad", Name: "Bad", BaseUrl: "http://x"},
			testLanguageCatalog(),
		)
		return ok
	})
	assertBadRequest(t, func(w http.ResponseWriter) bool {
		_, ok := metadataProviderInput(w, MetadataProviderRequest{Name: "Bad", Type: Tmdb, BaseUrl: "http://x", Priority: 1001})
		return ok
	})
}

func testLanguageCatalog() []storage.Language {
	return []storage.Language{
		{Code: "EN", DisplayName: "English", Aliases: []string{"ENG"}},
		{Code: "DE", DisplayName: "German", Aliases: []string{"Deutsch", "DEU"}},
	}
}

func TestSCNSettings009UserTagLanguageAndQualityInputsValidateEdges(t *testing.T) {
	user, ok := userCreateInput(httptest.NewRecorder(), UserCreateRequest{
		Username: " scenario ",
		Password: " password123 ",
		Role:     Admin,
	})
	if !ok || user.Username != "scenario" || user.Password == nil || *user.Password != "password123" {
		t.Fatalf("user create input = %#v, ok = %v", user, ok)
	}

	updated, ok := userUpdateInput(httptest.NewRecorder(), UserUpdateRequest{
		Username: "viewer",
		Password: stringPtr(""),
		Role:     User,
	})
	if !ok || updated.Password != nil {
		t.Fatalf("user update input = %#v, ok = %v", updated, ok)
	}

	tag, ok := tagInput(httptest.NewRecorder(), TagRequest{Name: "  Scenario   Tag  "})
	if !ok || tag != "Scenario Tag" {
		t.Fatalf("tag input = %q, ok = %v", tag, ok)
	}

	language, ok := languageInput(httptest.NewRecorder(), LanguageRequest{
		Code:        "en",
		DisplayName: " English ",
		Aliases:     []string{"eng"},
	})
	if !ok || language.DisplayName != " English " {
		t.Fatalf("language input = %#v, ok = %v", language, ok)
	}

	preferred := 2.0
	maximum := 3.0
	sizes, ok := qualitySizeSettingsInput(httptest.NewRecorder(), QualitySizeSettingsUpdateRequest{
		Qualities: []QualitySizeSettingRequest{{
			QualityId:                " webdl-1080p ",
			MinimumSizeMbPerMinute:   1,
			PreferredSizeMbPerMinute: &preferred,
			MaximumSizeMbPerMinute:   &maximum,
		}},
	})
	if !ok || len(sizes) != 1 || sizes[0].QualityID != "webdl-1080p" {
		t.Fatalf("quality size input = %#v, ok = %v", sizes, ok)
	}

	assertBadRequest(t, func(w http.ResponseWriter) bool {
		_, ok := userCreateInput(w, UserCreateRequest{Username: "x", Password: "short", Role: Admin})
		return ok
	})
	assertBadRequest(t, func(w http.ResponseWriter) bool {
		_, ok := tagInput(w, TagRequest{Name: "  "})
		return ok
	})
	assertBadRequest(t, func(w http.ResponseWriter) bool {
		_, ok := languageInput(w, LanguageRequest{Code: "en"})
		return ok
	})
	assertBadRequest(t, func(w http.ResponseWriter) bool {
		_, ok := qualitySizeSettingsInput(w, QualitySizeSettingsUpdateRequest{
			Qualities: []QualitySizeSettingRequest{{QualityId: "q", MinimumSizeMbPerMinute: -1}},
		})
		return ok
	})
}

func TestSCNSettings019FileNamingInputRequiresAllTemplates(t *testing.T) {
	input, ok := fileNamingSettingsInput(httptest.NewRecorder(), FileNamingSettingsRequest{
		MovieFileFormat:      " {movie_title} ",
		MovieFolderFormat:    "{movie_title}",
		SeriesEpisodeFormat:  "{series_title}",
		DailyEpisodeFormat:   "{series_title} {air_date}",
		AnimeEpisodeFormat:   "{series_title} {absolute_episode}",
		SeriesFolderFormat:   "{series_title}",
		SeasonFolderFormat:   "Season {season}",
		SpecialsFolderFormat: "Specials",
	})
	if !ok || input.MovieFileFormat != "{movie_title}" {
		t.Fatalf("file naming input = %#v, ok = %v", input, ok)
	}

	assertBadRequest(t, func(w http.ResponseWriter) bool {
		_, ok := fileNamingSettingsInput(w, FileNamingSettingsRequest{MovieFileFormat: "{movie_title}"})
		return ok
	})
}

func assertBadRequest(t *testing.T, call func(http.ResponseWriter) bool) {
	t.Helper()
	response := httptest.NewRecorder()
	if call(response) {
		t.Fatal("expected mapper to reject request")
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %q", response.Code, http.StatusBadRequest, response.Body.String())
	}
}
