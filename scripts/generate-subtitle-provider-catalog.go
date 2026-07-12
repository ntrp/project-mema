//go:build subtitlecatalog

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"media-manager/internal/subtitles/catalog"
)

const subtitleCatalogOutPath = "internal/subtitles/catalog/subtitle_provider_catalog.generated.json"

type subtitleProviderDefinition = catalog.Entry

func main() {
	entries := append([]catalog.Entry(nil), subtitleProviderCatalogDefinitions...)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
	if err := validateSubtitleCatalog(entries); err != nil {
		fmt.Fprintf(os.Stderr, "subtitle catalog validation failed: %v\n", err)
		os.Exit(1)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode subtitle catalog: %v\n", err)
		os.Exit(1)
	}
	data = append(bytes.TrimSpace(data), '\n')
	if err := os.WriteFile(subtitleCatalogOutPath, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write subtitle catalog: %v\n", err)
		os.Exit(1)
	}
}

func validateSubtitleCatalog(entries []catalog.Entry) error {
	seen := map[string]struct{}{}
	for _, entry := range entries {
		if entry.Key == "" {
			return fmt.Errorf("empty provider key")
		}
		if _, ok := seen[entry.Key]; ok {
			return fmt.Errorf("duplicate provider key %q", entry.Key)
		}
		seen[entry.Key] = struct{}{}
		if entry.DisplayName == "" {
			return fmt.Errorf("provider %q has empty display name", entry.Key)
		}
		if entry.Key != "mock" && entry.ProvenanceCommit != bazarrProviderListCommit {
			return fmt.Errorf("provider %q has invalid provenance", entry.Key)
		}
		if entry.RuntimeStatus == "" || entry.RuntimeMessage == "" {
			return fmt.Errorf("provider %q has incomplete runtime state", entry.Key)
		}
		for _, field := range entry.Fields {
			if field.Key == "" || field.Label == "" || field.Type == "" {
				return fmt.Errorf("provider %q has incomplete field", entry.Key)
			}
			if field.Type == catalog.FieldAction && field.Persisted {
				return fmt.Errorf("provider %q action field %q must not persist", entry.Key, field.Key)
			}
		}
	}
	if _, ok := seen["opensubtitles"]; ok {
		return fmt.Errorf("legacy opensubtitles alias must not be a picker entry")
	}
	if len(entries) != 59 {
		return fmt.Errorf("expected 59 picker entries, got %d", len(entries))
	}
	return nil
}

func baseURLField(defaultValue string) catalog.Field {
	field := catalog.Field{Key: "baseUrl", Label: "Base URL", Type: catalog.FieldText, Required: true, Persisted: true, SemanticKey: "base_url"}
	if defaultValue != "" {
		field.Options = []string{defaultValue}
	}
	return field
}

func usernameField() catalog.Field {
	return catalog.Field{Key: "username", Label: "Username", Type: catalog.FieldText, Persisted: true, SemanticKey: "username"}
}

func passwordField() catalog.Field {
	return catalog.Field{Key: "password", Label: "Password", Type: catalog.FieldPassword, Secret: true, Persisted: true, SemanticKey: "password"}
}

func apiKeyField() catalog.Field {
	return catalog.Field{Key: "apiKey", Label: "API key", Type: catalog.FieldPassword, Secret: true, Required: true, Persisted: true, SemanticKey: "api_key"}
}
