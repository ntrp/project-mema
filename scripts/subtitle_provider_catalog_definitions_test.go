//go:build subtitlecatalog

package main

import (
	"testing"

	"media-manager/internal/subtitles/catalog"
)

func TestSubtitleProviderDefinitionsAreUniqueAndComplete(t *testing.T) {
	if err := validateSubtitleCatalog(subtitleProviderCatalogDefinitions); err != nil {
		t.Fatalf("catalog definitions should validate: %v", err)
	}
	seen := map[string]bool{}
	for _, entry := range subtitleProviderCatalogDefinitions {
		if seen[entry.Key] {
			t.Fatalf("duplicate key %q", entry.Key)
		}
		seen[entry.Key] = true
	}
	if len(seen) != 59 {
		t.Fatalf("expected 59 unique definitions, got %d", len(seen))
	}
}

func TestCatalogFieldHelpers(t *testing.T) {
	assertFieldHelper(t, cookiesField(), "cookies", catalog.FieldPassword, true, false, "cookies")
	assertFieldHelper(t, userAgentField(), "userAgent", catalog.FieldText, false, false, "user_agent")
	assertFieldHelper(t, tokenField(), "token", catalog.FieldPassword, true, true, "token")
	assertFieldHelper(t, passkeyField(), "passkey", catalog.FieldPassword, true, true, "passkey")
	assertFieldHelper(t, hashedPasswordField(), "hashedPassword", catalog.FieldPassword, true, true, "hashed_password")
	assertFieldHelper(t, boolField("vip", "VIP"), "vip", catalog.FieldSwitch, false, false, "")
	assertFieldHelper(t, numericTextField("timeout", "Timeout"), "timeout", catalog.FieldText, false, false, "")
}

func assertFieldHelper(t *testing.T, field catalog.Field, key string, fieldType catalog.FieldType, secret bool, required bool, semantic string) {
	t.Helper()
	if field.Key != key || field.Type != fieldType || field.Secret != secret || field.Required != required || !field.Persisted || field.SemanticKey != semantic {
		t.Fatalf("unexpected field helper result: %#v", field)
	}
}
