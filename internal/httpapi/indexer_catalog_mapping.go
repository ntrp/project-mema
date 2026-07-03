package httpapi

import (
	"encoding/json"

	"media-manager/internal/indexers"
)

func indexerCatalogResponse(entries []indexers.CatalogEntry) IndexerCatalogResponse {
	return IndexerCatalogResponse{
		Entries:    catalogEntryResponses(entries),
		Protocols:  catalogProtocols(entries),
		Languages:  catalogLanguages(entries),
		Privacy:    catalogPrivacy(entries),
		Categories: catalogCategories(entries),
	}
}

func catalogEntryResponses(entries []indexers.CatalogEntry) []IndexerCatalogEntry {
	responses := make([]IndexerCatalogEntry, 0, len(entries))
	for _, entry := range entries {
		responses = append(responses, catalogEntryResponse(entry))
	}
	return responses
}

func catalogEntryResponse(entry indexers.CatalogEntry) IndexerCatalogEntry {
	indexerURLs := append([]string(nil), entry.IndexerURLs...)
	legacyURLs := append([]string(nil), entry.LegacyURLs...)
	return IndexerCatalogEntry{
		DefinitionId:       entry.DefinitionID,
		Name:               entry.Name,
		Implementation:     entry.Implementation,
		ImplementationName: entry.ImplementationName,
		Description:        optionalCatalogString(entry.Description),
		Language:           entry.Language,
		Encoding:           optionalCatalogString(entry.Encoding),
		IndexerUrls:        &indexerURLs,
		LegacyUrls:         &legacyURLs,
		Protocol:           IndexerProtocol(entry.Protocol),
		Privacy:            IndexerPrivacy(entry.Privacy),
		SupportsRss:        entry.SupportsRSS,
		SupportsSearch:     entry.SupportsSearch,
		SupportsRedirect:   entry.SupportsRedirect,
		SupportsPagination: entry.SupportsPagination,
		Capabilities:       catalogCapabilities(entry.Capabilities),
		Fields:             catalogFields(entry.Fields),
	}
}

func catalogCapabilities(capabilities indexers.Capabilities) IndexerCapabilities {
	return IndexerCapabilities{
		LimitsMax:         capabilities.LimitsMax,
		LimitsDefault:     capabilities.LimitsDefault,
		Categories:        catalogCategoryResponses(capabilities.Categories),
		SupportsRawSearch: capabilities.SupportsRawSearch,
		SearchParams:      append([]string(nil), capabilities.SearchParams...),
		TvSearchParams:    append([]string(nil), capabilities.TvSearchParams...),
		MovieSearchParams: append([]string(nil), capabilities.MovieSearchParams...),
	}
}

func catalogFields(fields []indexers.Field) []IndexerField {
	responses := make([]IndexerField, 0, len(fields))
	for _, field := range fields {
		responses = append(responses, IndexerField{
			Order:           &field.Order,
			Name:            field.Name,
			Label:           field.Label,
			Unit:            optionalCatalogString(field.Unit),
			HelpText:        optionalCatalogString(field.HelpText),
			HelpTextWarning: optionalCatalogString(field.HelpWarning),
			HelpLink:        optionalCatalogString(field.HelpLink),
			Value:           &field.Value,
			Type:            IndexerFieldType(field.Type),
			Advanced:        field.Advanced,
			SelectOptions:   catalogSelectOptions(field.SelectOptions),
			Section:         optionalCatalogString(field.Section),
			Placeholder:     optionalCatalogString(field.Placeholder),
			IsFloat:         &field.IsFloat,
		})
	}
	return responses
}

func catalogCategoryResponses(categories []indexers.Category) []IndexerCategory {
	responses := make([]IndexerCategory, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, IndexerCategory{
			Id:       category.ID,
			Name:     category.Name,
			Children: catalogCategoryResponses(category.Children),
		})
	}
	return responses
}

func catalogSelectOptions(options []indexers.SelectOption) *[]IndexerFieldSelectOption {
	if len(options) == 0 {
		return nil
	}
	responses := make([]IndexerFieldSelectOption, 0, len(options))
	for _, option := range options {
		responses = append(responses, IndexerFieldSelectOption{Value: option.Value, Name: option.Name})
	}
	return &responses
}

func indexerFieldValues(raw json.RawMessage) *[]IndexerFieldValue {
	if len(raw) == 0 {
		return nil
	}
	values := []IndexerFieldValue{}
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil
	}
	return &values
}

func indexerCapabilities(raw json.RawMessage) IndexerCapabilities {
	var capabilities IndexerCapabilities
	if err := json.Unmarshal(raw, &capabilities); err != nil {
		return IndexerCapabilities{
			Categories:        []IndexerCategory{},
			SupportsRawSearch: true,
			SearchParams:      []string{"q"},
			TvSearchParams:    []string{"q", "season", "ep"},
			MovieSearchParams: []string{"q", "imdbid"},
		}
	}
	return capabilities
}

func catalogProtocols(entries []indexers.CatalogEntry) []IndexerProtocol {
	seen := map[string]bool{}
	values := []IndexerProtocol{}
	for _, entry := range entries {
		if seen[entry.Protocol] {
			continue
		}
		seen[entry.Protocol] = true
		values = append(values, IndexerProtocol(entry.Protocol))
	}
	return values
}

func catalogLanguages(entries []indexers.CatalogEntry) []string {
	seen := map[string]bool{}
	values := []string{}
	for _, entry := range entries {
		if seen[entry.Language] {
			continue
		}
		seen[entry.Language] = true
		values = append(values, entry.Language)
	}
	return values
}

func catalogPrivacy(entries []indexers.CatalogEntry) []IndexerPrivacy {
	seen := map[string]bool{}
	values := []IndexerPrivacy{}
	for _, entry := range entries {
		group := entry.Privacy
		if group == "semiPrivate" {
			group = "private"
		}
		if seen[group] {
			continue
		}
		seen[group] = true
		values = append(values, IndexerPrivacy(group))
	}
	return values
}

func catalogCategories(entries []indexers.CatalogEntry) []IndexerCategory {
	seen := map[int32]IndexerCategory{}
	for _, entry := range entries {
		for _, category := range catalogCategoryResponses(entry.Capabilities.Categories) {
			seen[category.Id] = category
		}
	}
	values := make([]IndexerCategory, 0, len(seen))
	for _, category := range seen {
		values = append(values, category)
	}
	return values
}
