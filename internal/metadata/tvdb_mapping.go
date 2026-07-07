package metadata

import (
	"strconv"
	"strings"
)

func tvdbDetailsResult(item tvdbDetails, mediaType string, externalID string) Details {
	date := tvdbPrimaryDate(item)
	details := Details{
		Title:            tvdbTitle(item),
		Type:             mediaType,
		Year:             tvdbYear(item, date),
		ExternalProvider: "tvdb",
		ExternalID:       externalID,
		ExternalURL:      optionalString(tvdbPageURL(mediaType, item.Slug, externalID)),
		Overview:         optionalString(tvdbOverview(item)),
		PosterPath:       optionalString(firstNonEmpty(item.Image, tvdbPoster(item.Artworks))),
		BackdropPath:     optionalString(tvdbBackdrop(item.Artworks)),
		TrailerURL:       optionalString(tvdbTrailerURL(item.Trailers)),
		Status:           optionalString(item.Status.String()),
		OriginalLanguage: optionalString(item.OriginalLanguage),
		RuntimeMinutes:   tvdbRuntime(item),
		Genres:           tvdbNames(item.Genres),
		Keywords:         tvdbKeywords(item),
		Facts:            []Fact{},
		Seasons:          []Season{},
		Cast:             []Person{},
		Crew:             []Person{},
	}
	if mediaType == "movie" {
		details.ReleaseDate = optionalString(date)
	} else {
		details.FirstAirDate = optionalString(date)
		details.Seasons = tvdbSeasons(item)
		if count := int32(len(details.Seasons)); count > 0 {
			details.SeasonCount = &count
		}
		if count := tvdbEpisodeCount(item, details.Seasons); count > 0 {
			details.EpisodeCount = &count
		}
	}
	details.Cast, details.Crew = tvdbPeople(item.Characters)
	details.Facts = append(details.Facts, tvdbCertificationFacts(item.ContentRatings)...)
	details.Facts = append(details.Facts, tvdbReleaseFacts(item.Releases)...)
	details.Facts = append(details.Facts, tvdbMoneyFacts(item)...)
	details.Facts = append(details.Facts, tvdbStringFact("Original Country", []string{tvdbCountryDisplay(item.OriginalCountry, "")})...)
	details.Facts = append(details.Facts, tvdbStringFact("Production Countries", tvdbProductionCountries(item.ProductionCountries))...)
	details.Facts = append(details.Facts, tvdbStringFact("Studios", tvdbCompanyNames(item))...)
	details.Facts = append(details.Facts, tvdbStringFact("Networks", tvdbNames(item.Companies.Network))...)
	details.Facts = append(details.Facts, tvdbRemoteIDFacts(item.RemoteIDs)...)
	details.Facts = append(details.Facts, tvdbStringFact("Spoken Languages", item.SpokenLanguages)...)
	details.Facts = append(details.Facts, tvdbStringFact("Subtitle Languages", item.SubtitleLanguages)...)
	return details
}

func tvdbPageURL(mediaType string, slug string, externalID string) string {
	slug = strings.Trim(strings.TrimSpace(slug), "/")
	switch mediaType {
	case "movie":
		if slug == "" {
			return tvdbDereferrerURL("movie", externalID)
		}
		return "https://thetvdb.com/movies/" + slug
	case "serie":
		if slug == "" {
			return tvdbDereferrerURL("series", externalID)
		}
		return "https://thetvdb.com/series/" + slug
	default:
		return ""
	}
}

func tvdbDereferrerURL(kind string, externalID string) string {
	externalID = strings.TrimSpace(externalID)
	if externalID == "" {
		return ""
	}
	return "https://thetvdb.com/dereferrer/" + kind + "/" + externalID
}

func tvdbTitle(item tvdbDetails) string {
	return firstNonEmpty(item.Name, item.Title)
}

func tvdbYear(item tvdbDetails, date string) *int32 {
	if year := yearFromString(item.Year.String()); year != nil {
		return year
	}
	return yearFromDate(date)
}

func tvdbPrimaryDate(item tvdbDetails) string {
	if item.FirstAired != "" {
		return item.FirstAired
	}
	if item.FirstRelease.String() != "" {
		return item.FirstRelease.String()
	}
	for _, release := range item.Releases {
		if strings.TrimSpace(release.Date) == "" {
			continue
		}
		if strings.EqualFold(release.Country, "usa") || strings.EqualFold(release.Country, "us") {
			return release.Date
		}
	}
	if len(item.Releases) > 0 {
		return item.Releases[0].Date
	}
	return ""
}

func tvdbRuntime(item tvdbDetails) *int32 {
	if item.Runtime > 0 {
		return &item.Runtime
	}
	if item.AverageRuntime > 0 {
		return &item.AverageRuntime
	}
	return nil
}

func tvdbNames(items []tvdbNamedEntity) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		if name := strings.TrimSpace(item.Name); name != "" {
			names = append(names, name)
		}
	}
	return names
}

func tvdbOverview(item tvdbDetails) string {
	if value := strings.TrimSpace(item.Overview); value != "" {
		return value
	}
	if value := tvdbPreferredTranslation(item.Translations.OverviewTranslations).Overview; value != "" {
		return value
	}
	return ""
}

func tvdbSeasons(item tvdbDetails) []Season {
	seasons := make([]Season, 0, len(item.Seasons))
	for _, season := range item.Seasons {
		name := strings.TrimSpace(season.Name)
		if name == "" && season.Number > 0 {
			name = "Season " + strconv.Itoa(int(season.Number))
		}
		if name == "" {
			continue
		}
		mapped := Season{
			Name:         name,
			SeasonNumber: season.Number,
			PosterPath:   optionalString(season.Image),
			Episodes:     []Episode{},
		}
		if season.EpisodeCount > 0 {
			mapped.EpisodeCount = &season.EpisodeCount
		}
		seasons = append(seasons, mapped)
	}
	return seasons
}

func tvdbEpisodeCount(item tvdbDetails, seasons []Season) int32 {
	if len(item.Episodes) > 0 {
		return int32(len(item.Episodes))
	}
	var count int32
	for _, season := range seasons {
		if season.EpisodeCount != nil {
			count += *season.EpisodeCount
		}
	}
	return count
}

func tvdbPeople(characters []tvdbCharacter) ([]Person, []Person) {
	cast := []Person{}
	crew := []Person{}
	for _, character := range characters {
		person := tvdbPerson(character)
		if person.Name == "" {
			continue
		}
		if tvdbCharacterIsCast(character) {
			cast = append(cast, person)
			continue
		}
		crew = append(crew, person)
	}
	return cast, crew
}

func tvdbPerson(character tvdbCharacter) Person {
	externalID := firstNonEmpty(character.PeopleID.String(), character.PersonID.String())
	return Person{
		ExternalProvider: optionalString("tvdb"),
		ExternalID:       optionalString(externalID),
		Name:             firstNonEmpty(character.PersonName, character.Name),
		Role:             optionalString(firstNonEmpty(character.Role, character.Name, character.TypeName)),
		ProfilePath:      optionalString(character.Image),
	}
}

func tvdbCharacterIsCast(character tvdbCharacter) bool {
	kind := strings.ToLower(firstNonEmpty(character.PeopleType, character.TypeName))
	return kind == "" || strings.Contains(kind, "actor") || strings.Contains(kind, "cast") || character.Type == 3
}
