package metadata

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

func anilistAliases(media anilistMedia) []Alias {
	provider := "anilist"
	language := nativeLanguage(media.CountryOfOrigin)
	romanizedLanguage := optionalString("")
	if language != nil {
		romanizedLanguage = optionalString(*language + "-Latn")
	}
	values := []struct {
		value, kind, field string
		language           *string
	}{
		{media.Title.Romaji, "romaji", "title.romaji", romanizedLanguage},
		{media.Title.English, "english", "title.english", optionalString("en")},
		{media.Title.Native, "native", "title.native", language},
	}
	for _, synonym := range media.Synonyms {
		values = append(values, struct {
			value, kind, field string
			language           *string
		}{synonym, "synonym", "synonyms", nil})
	}
	aliases := make([]Alias, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		trimmed := strings.TrimSpace(value.value)
		key := strings.ToLower(trimmed) + "\x00" + value.kind
		if trimmed == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		aliases = append(aliases, Alias{Value: trimmed, Language: value.language, Kind: value.kind, ProviderName: provider, Source: map[string]any{
			"source": "anilist", "providerId": strconv.Itoa(media.ID), "field": value.field,
		}})
	}
	return aliases
}

func nativeLanguage(country string) *string {
	languages := map[string]string{"JP": "ja", "CN": "zh", "KR": "ko", "TW": "zh"}
	return optionalString(languages[strings.ToUpper(strings.TrimSpace(country))])
}

func anilistMappings(media anilistMedia, mediaType string) []ProviderMapping {
	raw := map[string]any{}
	payload := media.Raw
	if len(payload) == 0 {
		payload, _ = json.Marshal(media)
	}
	_ = json.Unmarshal(payload, &raw)
	mappings := []ProviderMapping{{
		ProviderName:       "anilist",
		ProviderEntityType: mediaType,
		ExternalID:         strconv.Itoa(media.ID),
		Canonical:          true,
		Source:             map[string]any{"source": "anilist", "raw": raw},
	}}
	seen := map[string]struct{}{"anilist:" + strconv.Itoa(media.ID): {}}
	for _, link := range media.ExternalLinks {
		provider, entityType, externalID := supportedExternalLink(link.Site, link.URL, mediaType)
		key := provider + ":" + externalID
		if provider == "" || externalID == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		mappings = append(mappings, ProviderMapping{
			ProviderName: provider, ProviderEntityType: entityType, ExternalID: externalID,
			Source: map[string]any{"source": "anilist_external_link", "site": link.Site, "url": link.URL},
		})
	}
	return mappings
}

func supportedExternalLink(_ string, rawURL, fallbackType string) (string, string, string) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", "", ""
	}
	provider := ""
	switch strings.ToLower(parsed.Hostname()) {
	case "anidb.net", "www.anidb.net":
		provider = "anidb"
	case "themoviedb.org", "www.themoviedb.org":
		provider = "tmdb"
	case "thetvdb.com", "www.thetvdb.com":
		provider = "tvdb"
	default:
		return "", "", ""
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	externalID := ""
	for index := len(parts) - 1; index >= 0; index-- {
		if _, err := strconv.Atoi(parts[index]); err == nil {
			externalID = parts[index]
			break
		}
	}
	entityType := fallbackType
	if provider == "tmdb" && len(parts) > 0 {
		if parts[0] == "movie" {
			entityType = "movie"
		} else if parts[0] == "tv" {
			entityType = "serie"
		}
	}
	return provider, entityType, externalID
}
