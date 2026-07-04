package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

const newznabSourcePath = "src/NzbDrone.Core/Indexers/Definitions/Newznab/Newznab.cs"

func fetchNewznabDefaults(client *http.Client) ([]catalogEntry, error) {
	req, err := http.NewRequest(http.MethodGet, prowlarrSourceRawBaseURL()+newznabSourcePath, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github returned %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return newznabEntriesFromSource(string(body)), nil
}

func prowlarrSourceRawBaseURL() string {
	if value := os.Getenv("PROWLARR_SOURCE_RAW_BASE_URL"); value != "" {
		return strings.TrimRight(value, "/") + "/"
	}
	return "https://raw.githubusercontent.com/Prowlarr/Prowlarr/develop/"
}

func newznabEntriesFromSource(source string) []catalogEntry {
	re := regexp.MustCompile(`(?s)yield\s+return\s+GetDefinition\("([^"]+)",\s*GetSettings\("([^"]*)"(?:,\s*apiPath:\s*@?"([^"]+)")?\)(?:,\s*categories:\s*new\[\]\s*\{([^}]*)\})?\);`)
	matches := re.FindAllStringSubmatch(source, -1)
	entries := make([]catalogEntry, 0, len(matches))
	for _, match := range matches {
		name := match[1]
		baseURL := match[2]
		if baseURL == "" || strings.EqualFold(name, "Generic Newznab") {
			continue
		}
		apiPath := match[3]
		if apiPath == "" {
			apiPath = "/api"
		}
		entries = append(entries, newznabCatalogEntry(name, baseURL, apiPath, categoryIDs(match[4])))
	}
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})
	return entries
}

func newznabCatalogEntry(name string, baseURL string, apiPath string, categoryIDs []int32) catalogEntry {
	limits := int32(100)
	return catalogEntry{
		DefinitionID:       "newznab-" + slugID(name),
		Name:               name,
		Implementation:     "Newznab",
		ImplementationName: name,
		Protocol:           "usenet",
		Privacy:            "private",
		Language:           "en-US",
		Description:        name + " Usenet indexer via Newznab API",
		IndexerURLs:        []string{joinURLPath(baseURL, apiPath)},
		LegacyURLs:         []string{},
		SupportsRSS:        true,
		SupportsSearch:     true,
		SupportsRedirect:   true,
		SupportsPagination: true,
		Capabilities: capabilities{
			LimitsMax:         &limits,
			LimitsDefault:     &limits,
			Categories:        newznabCategories(categoryIDs),
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

func categoryIDs(raw string) []int32 {
	values := []int32{}
	for _, item := range strings.Split(raw, ",") {
		parsed, err := strconv.ParseInt(strings.TrimSpace(item), 10, 32)
		if err == nil {
			values = append(values, int32(parsed))
		}
	}
	return values
}

func newznabCategories(ids []int32) []category {
	seen := map[int32]string{}
	for _, id := range ids {
		top := (id / 1000) * 1000
		if name, ok := categoryName(top); ok {
			seen[top] = name
		}
	}
	if len(seen) == 0 {
		for _, id := range []int32{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000} {
			if name, ok := categoryName(id); ok {
				seen[id] = name
			}
		}
	}
	cats := make([]category, 0, len(seen))
	for id, name := range seen {
		cats = append(cats, category{ID: id, Name: name, Children: []category{}})
	}
	sort.Slice(cats, func(i, j int) bool { return cats[i].ID < cats[j].ID })
	return cats
}

func categoryName(id int32) (string, bool) {
	values := map[int32]string{1000: "Console", 2000: "Movies", 3000: "Audio", 4000: "PC", 5000: "TV", 6000: "XXX", 7000: "Books", 8000: "Other"}
	name, ok := values[id]
	return name, ok
}

func joinURLPath(baseURL string, path string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func slugID(value string) string {
	var out strings.Builder
	lastHyphen := false
	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			out.WriteRune(r)
			lastHyphen = false
			continue
		}
		if !lastHyphen {
			out.WriteByte('-')
			lastHyphen = true
		}
	}
	return strings.Trim(out.String(), "-")
}
