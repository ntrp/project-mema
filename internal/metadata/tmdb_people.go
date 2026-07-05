package metadata

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

func (s *Service) PersonDetails(ctx context.Context, config Config, personID string) (PersonDetails, error) {
	if config.Type != "tmdb" {
		return PersonDetails{}, ErrUnsupportedProvider
	}
	personID = strings.TrimSpace(personID)
	if personID == "" {
		return PersonDetails{}, ErrUnsupportedProvider
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "person", personID)
	if err != nil {
		return PersonDetails{}, err
	}
	values := url.Values{}
	values.Set("append_to_response", "combined_credits")
	endpoint = endpoint + "?" + values.Encode()

	var payload tmdbPersonDetails
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return PersonDetails{}, err
	}
	return tmdbPersonDetailsResult(payload, personID), nil
}

func tmdbPersonDetailsResult(item tmdbPersonDetails, personID string) PersonDetails {
	details := PersonDetails{
		ID:           personID,
		Name:         strings.TrimSpace(item.Name),
		Biography:    optionalString(item.Biography),
		Birthday:     optionalString(item.Birthday),
		Deathday:     optionalString(item.Deathday),
		PlaceOfBirth: optionalString(item.PlaceOfBirth),
		ProfilePath:  optionalString(item.ProfilePath),
		AlsoKnownAs:  cleanedStrings(item.AlsoKnownAs),
		Appearances:  combinedAppearances(item.Credits),
	}
	return details
}

func combinedAppearances(credits tmdbCombinedCredits) []PersonAppearance {
	seen := map[string]int{}
	appearances := []PersonAppearance{}
	for _, item := range credits.Cast {
		appendAppearance(&appearances, seen, tmdbCreditAppearance(item, item.Character))
	}
	for _, item := range credits.Crew {
		appendAppearance(&appearances, seen, tmdbCreditAppearance(item, item.Job))
	}
	slices.SortFunc(appearances, func(a PersonAppearance, b PersonAppearance) int {
		return strings.Compare(stringValue(b.ReleaseDate), stringValue(a.ReleaseDate))
	})
	return appearances
}

func appendAppearance(items *[]PersonAppearance, seen map[string]int, item PersonAppearance) {
	if item.Title == "" || item.Type == "" {
		return
	}
	key := item.Type + ":" + item.ExternalID
	if index, ok := seen[key]; ok {
		merged := (*items)[index]
		merged.Role = mergeRoles(merged.Role, item.Role)
		(*items)[index] = merged
		return
	}
	seen[key] = len(*items)
	*items = append(*items, item)
}

func tmdbCreditAppearance(item tmdbCreditMedia, role string) PersonAppearance {
	mediaType := tmdbResultMediaType(item.MediaType)
	title := strings.TrimSpace(item.Title)
	date := item.ReleaseDate
	if mediaType == "serie" {
		title = strings.TrimSpace(item.Name)
		date = item.FirstAirDate
	}
	return PersonAppearance{
		Title:            title,
		Type:             mediaType,
		Year:             yearFromDate(date),
		ExternalProvider: "tmdb",
		ExternalID:       strconv.FormatInt(item.ID, 10),
		Overview:         optionalString(item.Overview),
		PosterPath:       optionalString(item.PosterPath),
		BackdropPath:     optionalString(item.BackdropPath),
		Role:             optionalString(role),
		ReleaseDate:      optionalString(date),
	}
}

func mergeRoles(current *string, next *string) *string {
	if current == nil {
		return next
	}
	if next == nil || strings.Contains(*current, *next) {
		return current
	}
	value := *current + ", " + *next
	return &value
}

func cleanedStrings(values []string) []string {
	items := []string{}
	for _, value := range values {
		if cleaned := strings.TrimSpace(value); cleaned != "" {
			items = append(items, cleaned)
		}
	}
	return items
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
