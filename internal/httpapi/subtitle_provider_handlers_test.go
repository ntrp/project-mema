package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"media-manager/internal/storage"
)

func TestSubtitleProviderResponseRedactsSecrets(t *testing.T) {
	response := subtitleProviderResponse(storage.SubtitleProvider{
		Name:            "OpenSubtitles",
		Type:            "opensubtitles",
		BaseURL:         "https://api.opensubtitles.com",
		Username:        stringPtr("user"),
		Password:        stringPtr("secret"),
		APIKey:          stringPtr("key"),
		SecretFieldsSet: []string{"apiKey", "password"},
	})
	if !response.ApiKeySet || !response.PasswordSet {
		t.Fatalf("expected secret presence flags, got %#v", response)
	}
	if response.CatalogKey != "opensubtitlescom" {
		t.Fatalf("expected legacy alias catalog key, got %q", response.CatalogKey)
	}
	if len(response.SecretFieldsSet) != 2 {
		t.Fatalf("expected secret field names, got %#v", response.SecretFieldsSet)
	}
}

func TestSubtitleProviderInputValidation(t *testing.T) {
	baseURL := "https://api.opensubtitles.com"
	w := httptest.NewRecorder()
	_, ok := subtitleProviderInput(w, SubtitleProviderRequest{
		Name:     "OpenSubtitles",
		Type:     Opensubtitles,
		BaseUrl:  &baseURL,
		Enabled:  true,
		Priority: 1001,
	})
	if ok || w.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid priority, ok=%v code=%d", ok, w.Code)
	}
}

func TestSubtitleProviderInputRejectsEnabledUnsupportedProvider(t *testing.T) {
	w := httptest.NewRecorder()
	_, ok := subtitleProviderInput(w, SubtitleProviderRequest{
		Name:     "BSPlayer",
		Type:     Bsplayer,
		Enabled:  true,
		Priority: 10,
	})
	if ok || w.Code != http.StatusBadRequest {
		t.Fatalf("expected unsupported enabled provider to be rejected, ok=%v code=%d", ok, w.Code)
	}
}

func TestSubtitleProviderInputAllowsDisabledCatalogProvider(t *testing.T) {
	w := httptest.NewRecorder()
	input, ok := subtitleProviderInput(w, SubtitleProviderRequest{
		Name:     "BSPlayer",
		Type:     Bsplayer,
		Enabled:  false,
		Priority: 10,
	})
	if !ok || input.Enabled || w.Code != http.StatusOK {
		t.Fatalf("expected disabled catalog provider to be accepted, input=%#v ok=%v code=%d", input, ok, w.Code)
	}
}

func TestSubtitleProviderConfigIncludesDynamicSettingsAndSecrets(t *testing.T) {
	text := "enabled"
	secret := "cookie=value"
	providerConfig := subtitleProviderConfig(storage.SubtitleProvider{
		Name:           "Dynamic Provider",
		Type:           "mock",
		Settings:       storage.SubtitleProviderSettings{"mode": {StringValue: &text}},
		SecretSettings: storage.SubtitleProviderSecretSettings{"cookies": secret},
	})
	if providerConfig.Settings["mode"].StringValue == nil || *providerConfig.Settings["mode"].StringValue != "enabled" {
		t.Fatalf("settings = %#v", providerConfig.Settings)
	}
	if providerConfig.SecretSettings["cookies"] != secret {
		t.Fatalf("secrets = %#v", providerConfig.SecretSettings)
	}
}

func TestDraftSubtitleProviderConfigIncludesDynamicSettingsAndSecrets(t *testing.T) {
	text := "fallback"
	secret := "token"
	input := storage.SubtitleProviderInput{
		Name:           "Draft",
		Type:           "mock",
		Settings:       storage.SubtitleProviderSettings{"mode": {StringValue: &text}},
		SecretSettings: storage.SubtitleProviderSecretSettings{"apiToken": secret},
	}
	providerConfig := subtitleProviderConfigFromInput(input)
	if providerConfig.Settings["mode"].StringValue == nil || *providerConfig.Settings["mode"].StringValue != "fallback" {
		t.Fatalf("settings = %#v", providerConfig.Settings)
	}
	if providerConfig.SecretSettings["apiToken"] != secret {
		t.Fatalf("secrets = %#v", providerConfig.SecretSettings)
	}
}

func TestSubtitleProviderUpdatePreservesAndClearsSecrets(t *testing.T) {
	input := storage.SubtitleProviderInput{SecretSettings: storage.SubtitleProviderSecretSettings{}}
	current := storage.SubtitleProvider{
		APIKey:         stringPtr("key"),
		Password:       stringPtr("secret"),
		SecretSettings: storage.SubtitleProviderSecretSettings{"apiKey": "key", "password": "secret"},
	}
	request := SubtitleProviderRequest{}
	input = preserveSubtitleProviderSecrets(input, request, current)

	if input.APIKey == nil || *input.APIKey != "key" {
		t.Fatalf("expected api key to be preserved")
	}
	if input.Password == nil || *input.Password != "secret" {
		t.Fatalf("expected password to be preserved")
	}

	clear := []string{"apiKey"}
	input.ClearSecretFields = clear
	input = preserveSubtitleProviderSecrets(input, SubtitleProviderRequest{ClearSecretFields: &clear}, current)
	if input.APIKey != nil || input.SecretSettings["apiKey"] != "" {
		t.Fatalf("expected api key to be cleared, got %#v", input)
	}
}
