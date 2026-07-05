package indexers

import (
	_ "embed"
	"encoding/json"
	"sort"
	"strings"
)

//go:embed indexer_catalog.generated.json
var indexerCatalogJSON []byte

type CatalogEntry struct {
	DefinitionID       string
	Name               string
	Implementation     string
	ImplementationName string
	Protocol           string
	Privacy            string
	Language           string
	Encoding           string
	Description        string
	IndexerURLs        []string
	LegacyURLs         []string
	SupportsRSS        bool
	SupportsSearch     bool
	SupportsRedirect   bool
	SupportsPagination bool
	Capabilities       Capabilities
	Fields             []Field
}

type Capabilities struct {
	LimitsMax         *int32
	LimitsDefault     *int32
	Categories        []Category
	SupportsRawSearch bool
	SearchParams      []string
	TvSearchParams    []string
	MovieSearchParams []string
}

type Category struct {
	ID       int32
	Name     string
	Children []Category
}

type Field struct {
	Order         int32
	Name          string
	Label         string
	Unit          string
	HelpText      string
	HelpWarning   string
	HelpLink      string
	Value         any
	Type          string
	Advanced      bool
	SelectOptions []SelectOption
	Section       string
	Placeholder   string
	IsFloat       bool
}

type SelectOption struct {
	Value string
	Name  string
}

func Catalog() []CatalogEntry {
	return append([]CatalogEntry(nil), catalogEntries...)
}

func CatalogEntryByID(definitionID string) (CatalogEntry, bool) {
	for _, entry := range catalogEntries {
		if entry.DefinitionID == definitionID {
			return entry, true
		}
	}
	return CatalogEntry{}, false
}

func DefaultMediaTypeScopes(entry CatalogEntry) []string {
	seen := map[string]bool{}
	collectCategoryScopes(entry.Capabilities.Categories, seen)
	if len(seen) == 0 {
		return []string{"movie", "serie", "anime", "audio", "book"}
	}
	order := []string{"movie", "serie", "anime", "audio", "book"}
	scopes := []string{}
	for _, scope := range order {
		if seen[scope] {
			scopes = append(scopes, scope)
		}
	}
	return scopes
}

func collectCategoryScopes(categories []Category, seen map[string]bool) {
	for _, category := range categories {
		name := strings.ToLower(category.Name)
		switch {
		case strings.Contains(name, "anime"):
			seen["anime"] = true
		case strings.Contains(name, "movie"):
			seen["movie"] = true
		case strings.Contains(name, "tv") || strings.Contains(name, "series") || strings.Contains(name, "show"):
			seen["serie"] = true
		case strings.Contains(name, "book") || strings.Contains(name, "ebook") || strings.Contains(name, "audiobook"):
			seen["book"] = true
		case strings.Contains(name, "audio") || strings.Contains(name, "music"):
			seen["audio"] = true
		}
		collectCategoryScopes(category.Children, seen)
	}
}

var catalogEntries = loadCatalogEntries()

func loadCatalogEntries() []CatalogEntry {
	entries := []CatalogEntry{}
	if err := json.Unmarshal(indexerCatalogJSON, &entries); err != nil {
		panic("load indexer catalog: " + err.Error())
	}
	for index := range entries {
		normalizeCatalogEntry(&entries[index])
	}
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})
	return entries
}

func normalizeCatalogEntry(entry *CatalogEntry) {
	if entry.Language == "" {
		entry.Language = "en-US"
	}
	if entry.Privacy == "" {
		entry.Privacy = "private"
	}
	if entry.Protocol == "" {
		entry.Protocol = "torrent"
	}
	if entry.Implementation == "" {
		entry.Implementation = "Cardigann"
	}
	if entry.ImplementationName == "" {
		entry.ImplementationName = entry.Name
	}
	if entry.IndexerURLs == nil {
		entry.IndexerURLs = []string{}
	}
	if entry.LegacyURLs == nil {
		entry.LegacyURLs = []string{}
	}
	if entry.Capabilities.Categories == nil {
		entry.Capabilities.Categories = []Category{}
	}
	if entry.Capabilities.SearchParams == nil {
		entry.Capabilities.SearchParams = []string{}
	}
	if entry.Capabilities.TvSearchParams == nil {
		entry.Capabilities.TvSearchParams = []string{}
	}
	if entry.Capabilities.MovieSearchParams == nil {
		entry.Capabilities.MovieSearchParams = []string{}
	}
	if entry.Fields == nil {
		entry.Fields = []Field{}
	}
}
