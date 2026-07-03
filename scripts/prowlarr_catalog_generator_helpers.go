package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func fieldType(value string) string {
	switch strings.ToLower(value) {
	case "password":
		return "password"
	case "checkbox":
		return "checkbox"
	case "select":
		return "select"
	case "info", "info_flaresolverr":
		return "info"
	case "number":
		return "number"
	default:
		return "text"
	}
}

func selectOptions(value any) []selectOption {
	options := []selectOption{}
	for key, rawName := range mapValue(value) {
		options = append(options, selectOption{Value: key, Name: stringValue(rawName)})
	}
	sort.Slice(options, func(i, j int) bool { return options[i].Name < options[j].Name })
	return options
}

func genericTorznab() catalogEntry {
	return genericEntry("generic-torznab", "Generic Torznab", "torrent", "Cardigann", "Private torrent tracker via Torznab API")
}

func genericNewznab() catalogEntry {
	return genericEntry("generic-newznab", "Generic Newznab", "usenet", "Newznab", "Usenet indexer via Newznab API")
}

func genericEntry(id, name, protocol, implementation, description string) catalogEntry {
	limits := int32(100)
	return catalogEntry{
		DefinitionID:       id,
		Name:               name,
		Implementation:     implementation,
		ImplementationName: name,
		Protocol:           protocol,
		Privacy:            "private",
		Language:           "en-US",
		Description:        description,
		SupportsRSS:        true,
		SupportsSearch:     true,
		SupportsRedirect:   true,
		SupportsPagination: true,
		IndexerURLs:        []string{},
		LegacyURLs:         []string{},
		Capabilities: capabilities{
			LimitsMax:         &limits,
			LimitsDefault:     &limits,
			Categories:        []category{{ID: 2000, Name: "Movies", Children: []category{}}, {ID: 5000, Name: "TV", Children: []category{}}},
			SupportsRawSearch: true,
			SearchParams:      []string{"q"},
			TvSearchParams:    []string{"q", "season", "ep"},
			MovieSearchParams: []string{"q", "imdbid"},
		},
		Fields: []catalogField{
			{Order: 1, Name: "baseUrl", Label: "Base URL", Type: "url"},
			{Order: 2, Name: "apiKey", Label: "API key", Type: "password"},
			{Order: 3, Name: "categories", Label: "Categories", Type: "text"},
		},
	}
}

func getJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github returned %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func mapValue(value any) map[string]any {
	result, _ := value.(map[string]any)
	return result
}

func arrayValue(value any) []any {
	result, _ := value.([]any)
	return result
}

func stringSlice(value any) []string {
	values := []string{}
	for _, item := range arrayValue(value) {
		values = append(values, stringValue(item))
	}
	return values
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return ""
	}
}

func boolValue(value any) bool {
	typed, _ := value.(bool)
	return typed
}

func fallbackString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func privacyValue(value string) string {
	switch strings.ToLower(strings.ReplaceAll(value, "-", "")) {
	case "public":
		return "public"
	case "semiprivate":
		return "semiPrivate"
	default:
		return "private"
	}
}

func hasSearch(raw map[string]any) bool {
	return len(stringSlice(mapValue(mapValue(raw["caps"])["modes"])["search"])) > 0
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
