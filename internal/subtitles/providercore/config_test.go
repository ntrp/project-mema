package providercore

import (
	"testing"

	"media-manager/internal/subtitles"
)

func TestConfigViewReadsTypedSettingsAndSecrets(t *testing.T) {
	text := " value "
	number := float64(42)
	flag := true
	config := subtitles.Config{
		BaseURL: "https://example.test/root/",
		APIKey:  stringPtr(" legacy-api "),
		Settings: map[string]subtitles.SettingValue{
			"text":    {StringValue: &text},
			"number":  {NumberValue: &number},
			"flag":    {BooleanValue: &flag},
			"strings": {StringValues: []string{" one ", "", "two"}},
		},
		SecretSettings: map[string]string{"cookies": "sid=1", "token": " secret "},
	}
	view := NewConfig(config)
	if view.StringSetting("text") != "value" || view.IntSetting("number") != 42 || !view.BoolSetting("flag") {
		t.Fatalf("typed settings were not read correctly")
	}
	if got := view.StringsSetting("strings"); len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("strings = %#v", got)
	}
	if view.Secret("apiKey") != "legacy-api" || view.Secret("token") != "secret" || view.CookieString() != "sid=1" {
		t.Fatalf("secrets were not read correctly")
	}
	if secret, ok := view.RequiredSecret("token"); !ok || secret != "secret" {
		t.Fatalf("required secret = %q ok=%v", secret, ok)
	}
	if _, ok := view.RequiredSecret("missing"); ok {
		t.Fatal("missing required secret should report false")
	}
	if view.BaseURL("https://fallback.test") != "https://example.test/root" {
		t.Fatalf("base URL = %q", view.BaseURL(""))
	}
}

func stringPtr(value string) *string { return &value }
