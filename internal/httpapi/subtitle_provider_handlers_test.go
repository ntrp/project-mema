package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"media-manager/internal/storage"
)

func TestSubtitleProviderResponseIncludesSecrets(t *testing.T) {
	response := subtitleProviderResponse(storage.SubtitleProvider{
		Name:     "OpenSubtitles",
		Type:     "opensubtitles",
		BaseURL:  "https://api.opensubtitles.com",
		Username: stringPtr("user"),
		Password: stringPtr("secret"),
		APIKey:   stringPtr("key"),
	})
	if !response.ApiKeySet || !response.PasswordSet {
		t.Fatalf("expected secret presence flags, got %#v", response)
	}
	if response.ApiKey == nil || *response.ApiKey != "key" {
		t.Fatalf("expected api key in response, got %#v", response.ApiKey)
	}
	if response.Password == nil || *response.Password != "secret" {
		t.Fatalf("expected password in response, got %#v", response.Password)
	}
}

func TestSubtitleProviderInputValidation(t *testing.T) {
	w := httptest.NewRecorder()
	_, ok := subtitleProviderInput(w, SubtitleProviderRequest{
		Name:     "OpenSubtitles",
		Type:     Opensubtitles,
		BaseUrl:  "https://api.opensubtitles.com",
		Enabled:  true,
		Priority: 1001,
	})
	if ok || w.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid priority, ok=%v code=%d", ok, w.Code)
	}
}

func TestSubtitleProviderUpdatePreservesOmittedSecrets(t *testing.T) {
	input := storage.SubtitleProviderInput{APIKey: nil, Password: nil}
	current := storage.SubtitleProvider{APIKey: stringPtr("key"), Password: stringPtr("secret")}
	request := SubtitleProviderRequest{}
	input = preserveSubtitleProviderSecrets(input, request, current)

	if input.APIKey == nil || *input.APIKey != "key" {
		t.Fatalf("expected api key to be preserved")
	}
	if input.Password == nil || *input.Password != "secret" {
		t.Fatalf("expected password to be preserved")
	}
}
