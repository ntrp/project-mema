package storage

import "testing"

func TestSubtitleProviderSettingsJSONRoundTrip(t *testing.T) {
	baseURL := "https://api.opensubtitles.com"
	settings := SubtitleProviderSettings{"baseUrl": {StringValue: &baseURL}}
	secrets := SubtitleProviderSecretSettings{"apiKey": "secret"}

	if got := subtitleProviderSettingsFromJSON(subtitleProviderSettingsJSON(settings)); *got["baseUrl"].StringValue != baseURL {
		t.Fatalf("settings round trip = %#v", got)
	}
	if got := subtitleProviderSecretsFromJSON(subtitleProviderSecretsJSON(secrets)); got["apiKey"] != "secret" {
		t.Fatalf("secrets round trip = %#v", got)
	}
}

func TestPreserveSubtitleProviderUpdateSecretsAllowsExplicitClear(t *testing.T) {
	input := preserveSubtitleProviderUpdateSecrets(SubtitleProviderInput{}, SubtitleProvider{
		APIKey:         stringPtr("key"),
		Password:       stringPtr("secret"),
		SecretSettings: SubtitleProviderSecretSettings{"apiKey": "key", "password": "secret"},
	})
	if input.APIKey == nil || *input.APIKey != "key" || input.Password == nil || *input.Password != "secret" {
		t.Fatalf("expected omitted secrets to be preserved, got %#v", input)
	}

	input = preserveSubtitleProviderUpdateSecrets(SubtitleProviderInput{
		ClearSecretFields: []string{"apiKey", "password"},
	}, SubtitleProvider{
		APIKey:         stringPtr("key"),
		Password:       stringPtr("secret"),
		SecretSettings: SubtitleProviderSecretSettings{"apiKey": "key", "password": "secret"},
	})
	if input.APIKey != nil || input.Password != nil || len(input.SecretSettings) != 0 {
		t.Fatalf("expected explicit clear to remove secrets, got %#v", input)
	}
}

func TestNormalizeSubtitleProviderInputMirrorsLegacyFields(t *testing.T) {
	password := "secret"
	input := normalizedSubtitleProviderInput(SubtitleProviderInput{
		Type:           "opensubtitlescom",
		Settings:       SubtitleProviderSettings{"baseUrl": subtitleProviderSettingValue("https://api.opensubtitles.com")},
		Password:       &password,
		SecretSettings: SubtitleProviderSecretSettings{},
	})
	if input.BaseURL != "https://api.opensubtitles.com" {
		t.Fatalf("base url = %q", input.BaseURL)
	}
	if input.SecretSettings["password"] != password {
		t.Fatalf("secret settings = %#v", input.SecretSettings)
	}
}
