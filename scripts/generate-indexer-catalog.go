package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	outPath = "internal/indexers/indexer_catalog.generated.json"
)

var (
	treeURL = catalogTreeURL()
	rawBase = catalogRawBaseURL()
)

type catalogEntry struct {
	DefinitionID       string         `json:"definitionId"`
	Name               string         `json:"name"`
	Implementation     string         `json:"implementation"`
	ImplementationName string         `json:"implementationName"`
	Protocol           string         `json:"protocol"`
	Privacy            string         `json:"privacy"`
	Language           string         `json:"language"`
	Encoding           string         `json:"encoding,omitempty"`
	Description        string         `json:"description,omitempty"`
	IndexerURLs        []string       `json:"indexerUrls"`
	LegacyURLs         []string       `json:"legacyUrls"`
	SupportsRSS        bool           `json:"supportsRss"`
	SupportsSearch     bool           `json:"supportsSearch"`
	SupportsRedirect   bool           `json:"supportsRedirect"`
	SupportsPagination bool           `json:"supportsPagination"`
	Capabilities       capabilities   `json:"capabilities"`
	Fields             []catalogField `json:"fields"`
}

type capabilities struct {
	LimitsMax         *int32     `json:"limitsMax,omitempty"`
	LimitsDefault     *int32     `json:"limitsDefault,omitempty"`
	Categories        []category `json:"categories"`
	SupportsRawSearch bool       `json:"supportsRawSearch"`
	SearchParams      []string   `json:"searchParams"`
	TvSearchParams    []string   `json:"tvSearchParams"`
	MovieSearchParams []string   `json:"movieSearchParams"`
}

type category struct {
	ID       int32      `json:"id"`
	Name     string     `json:"name"`
	Children []category `json:"children"`
}

type catalogField struct {
	Order         int32          `json:"order"`
	Name          string         `json:"name"`
	Label         string         `json:"label"`
	Unit          string         `json:"unit,omitempty"`
	HelpText      string         `json:"helpText,omitempty"`
	HelpWarning   string         `json:"helpWarning,omitempty"`
	HelpLink      string         `json:"helpLink,omitempty"`
	Value         any            `json:"value,omitempty"`
	Type          string         `json:"type"`
	Advanced      bool           `json:"advanced"`
	SelectOptions []selectOption `json:"selectOptions,omitempty"`
	Section       string         `json:"section,omitempty"`
	Placeholder   string         `json:"placeholder,omitempty"`
	IsFloat       bool           `json:"isFloat,omitempty"`
}

type selectOption struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

func main() {
	paths, err := definitionPaths()
	if err != nil {
		fatal(err)
	}
	entries := make([]catalogEntry, 0, len(paths)+2)
	entries = append(entries, genericTorznab(), genericNewznab())
	client := &http.Client{Timeout: 20 * time.Second}
	for _, path := range paths {
		entry, err := fetchDefinition(client, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
			continue
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile(outPath, append(data, '\n'), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("wrote %d catalog entries to %s\n", len(entries), outPath)
}

func definitionPaths() ([]string, error) {
	var tree struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"tree"`
	}
	if err := getJSON(treeURL, &tree); err != nil {
		return nil, err
	}
	paths := []string{}
	for _, item := range tree.Tree {
		if item.Type == "blob" && strings.HasPrefix(item.Path, "definitions/v11/") && strings.HasSuffix(item.Path, ".yml") {
			paths = append(paths, item.Path)
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func catalogTreeURL() string {
	if value := os.Getenv("INDEXER_CATALOG_TREE_URL"); value != "" {
		return value
	}
	return fmt.Sprintf("https://api.github.com/repos/%s/git/trees/master?recursive=1", catalogRepo())
}

func catalogRawBaseURL() string {
	if value := os.Getenv("INDEXER_CATALOG_RAW_BASE_URL"); value != "" {
		return value
	}
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/master/", catalogRepo())
}

func catalogRepo() string {
	if value := os.Getenv("INDEXER_CATALOG_REPO"); value != "" {
		return value
	}
	return strings.Join([]string{"Prow", "larr/Indexers"}, "")
}

func fetchDefinition(client *http.Client, path string) (catalogEntry, error) {
	req, err := http.NewRequest(http.MethodGet, rawBase+path, nil)
	if err != nil {
		return catalogEntry{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return catalogEntry{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return catalogEntry{}, fmt.Errorf("github returned %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return catalogEntry{}, err
	}
	var raw map[string]any
	if err := yaml.Unmarshal(body, &raw); err != nil {
		return catalogEntry{}, err
	}
	return entryFromYAML(raw, path), nil
}

func entryFromYAML(raw map[string]any, path string) catalogEntry {
	id := stringValue(raw["id"])
	if id == "" {
		id = strings.TrimSuffix(strings.TrimPrefix(path, "definitions/v11/"), ".yml")
	}
	name := fallbackString(stringValue(raw["name"]), id)
	return catalogEntry{
		DefinitionID:       id,
		Name:               name,
		Implementation:     "Cardigann",
		ImplementationName: name,
		Protocol:           "torrent",
		Privacy:            privacyValue(stringValue(raw["type"])),
		Language:           fallbackString(stringValue(raw["language"]), "en-US"),
		Encoding:           stringValue(raw["encoding"]),
		Description:        stringValue(raw["description"]),
		IndexerURLs:        stringSlice(raw["links"]),
		LegacyURLs:         stringSlice(raw["legacylinks"]),
		SupportsRSS:        true,
		SupportsSearch:     hasSearch(raw),
		SupportsRedirect:   true,
		SupportsPagination: true,
		Capabilities:       capabilitiesFrom(raw),
		Fields:             fieldsFrom(raw),
	}
}

func capabilitiesFrom(raw map[string]any) capabilities {
	caps := mapValue(raw["caps"])
	modes := mapValue(caps["modes"])
	limits := int32(100)
	return capabilities{
		LimitsMax:         &limits,
		LimitsDefault:     &limits,
		Categories:        categoriesFrom(caps["categorymappings"]),
		SupportsRawSearch: boolValue(caps["allowrawsearch"]),
		SearchParams:      stringSlice(modes["search"]),
		TvSearchParams:    stringSlice(modes["tv-search"]),
		MovieSearchParams: stringSlice(modes["movie-search"]),
	}
}

func categoriesFrom(value any) []category {
	seen := map[int32]string{}
	for _, item := range arrayValue(value) {
		top := strings.Split(stringValue(mapValue(item)["cat"]), "/")[0]
		if id, ok := standardCategory(top); ok {
			seen[id] = top
		}
	}
	cats := make([]category, 0, len(seen))
	for id, name := range seen {
		cats = append(cats, category{ID: id, Name: name, Children: []category{}})
	}
	sort.Slice(cats, func(i, j int) bool { return cats[i].ID < cats[j].ID })
	return cats
}

func standardCategory(name string) (int32, bool) {
	values := map[string]int32{"Console": 1000, "Movies": 2000, "Audio": 3000, "PC": 4000, "TV": 5000, "XXX": 6000, "Books": 7000, "Other": 8000}
	id, ok := values[name]
	return id, ok
}

func fieldsFrom(raw map[string]any) []catalogField {
	items := arrayValue(raw["settings"])
	fields := make([]catalogField, 0, len(items))
	for i, item := range items {
		setting := mapValue(item)
		name := stringValue(setting["name"])
		if name == "" {
			continue
		}
		typ := fieldType(stringValue(setting["type"]))
		helpText := stringValue(setting["help"])
		if typ == "info" && helpText == "" {
			helpText = stringValue(setting["default"])
		}
		fields = append(fields, catalogField{
			Order:         int32(i + 1),
			Name:          name,
			Label:         fallbackString(stringValue(setting["label"]), name),
			HelpText:      helpText,
			Value:         setting["default"],
			Type:          typ,
			Advanced:      boolValue(setting["advanced"]),
			SelectOptions: selectOptions(setting["options"]),
		})
	}
	return fields
}
