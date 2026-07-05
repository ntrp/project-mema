package metadata

import (
	"strconv"
	"strings"
)

func tmdbDetailsResult(item tmdbDetails, mediaType string, externalID string) Details {
	title := strings.TrimSpace(item.Title)
	date := item.ReleaseDate
	if mediaType == "serie" {
		title = strings.TrimSpace(item.Name)
		date = item.FirstAirDate
	}
	details := Details{
		Title:            title,
		Type:             mediaType,
		Year:             yearFromDate(date),
		ExternalProvider: "tmdb",
		ExternalID:       externalID,
		Overview:         optionalString(item.Overview),
		PosterPath:       optionalString(item.PosterPath),
		CollectionID:     tmdbCollectionID(item.Collection),
		CollectionName:   tmdbCollectionName(item.Collection),
		BackdropPath:     optionalString(item.BackdropPath),
		TrailerURL:       tmdbTrailerURL(item.Videos),
		Status:           optionalString(item.Status),
		OriginalLanguage: optionalString(item.OriginalLanguage),
		Genres:           tmdbNames(item.Genres),
		Keywords:         tmdbKeywordNames(item.Keywords),
		Facts:            []Fact{},
		Seasons:          []Season{},
		Cast:             []Person{},
		Crew:             []Person{},
		Recommendations:  tmdbResults(item.Recommendations.Results, mediaType, 20),
		Similar:          tmdbResults(item.Similar.Results, mediaType, 20),
	}
	if mediaType == "movie" {
		details.ReleaseDate = optionalString(item.ReleaseDate)
		if item.Runtime > 0 {
			details.RuntimeMinutes = &item.Runtime
		}
	} else {
		details.FirstAirDate = optionalString(item.FirstAirDate)
		if item.NumberOfSeasons > 0 {
			details.SeasonCount = &item.NumberOfSeasons
		}
		if item.NumberOfEpisodes > 0 {
			details.EpisodeCount = &item.NumberOfEpisodes
		}
		if len(item.EpisodeRunTime) > 0 && item.EpisodeRunTime[0] > 0 {
			value := item.EpisodeRunTime[0]
			details.RuntimeMinutes = &value
		}
		for _, season := range item.Seasons {
			name := strings.TrimSpace(season.Name)
			if name == "" {
				continue
			}
			mapped := Season{
				Name:       name,
				AirDate:    optionalString(season.AirDate),
				PosterPath: optionalString(season.PosterPath),
				Episodes:   []Episode{},
			}
			if season.EpisodeCount > 0 {
				mapped.EpisodeCount = &season.EpisodeCount
			}
			for _, episode := range season.Episodes {
				episodeName := strings.TrimSpace(episode.Name)
				if episodeName == "" {
					continue
				}
				mapped.Episodes = append(mapped.Episodes, Episode{
					Name:          episodeName,
					EpisodeNumber: episode.EpisodeNumber,
					Overview:      optionalString(episode.Overview),
					AirDate:       optionalString(episode.AirDate),
					StillPath:     optionalString(episode.StillPath),
				})
			}
			details.Seasons = append(details.Seasons, mapped)
		}
	}
	if item.VoteAverage > 0 {
		details.VoteAverage = &item.VoteAverage
	}
	details.Facts = append(details.Facts, tmdbCertificationFacts(item)...)
	details.Facts = append(details.Facts, tmdbReleaseDateFacts(item.ReleaseDates)...)
	details.Facts = append(details.Facts, tmdbFinancialFacts(item)...)
	if len(item.Countries) > 0 {
		details.Facts = append(details.Facts, Fact{Label: "Production Countries", Value: strings.Join(tmdbCountries(item.Countries), "\n")})
	}
	if len(item.Production) > 0 {
		details.Facts = append(details.Facts, Fact{Label: "Studios", Value: strings.Join(tmdbNames(item.Production), "\n")})
	}
	details.Facts = append(details.Facts, tmdbCrewFacts(item.Credits.Crew)...)
	if len(item.CreatedBy) > 0 {
		details.Facts = append(details.Facts, Fact{Label: "Creator", Value: strings.Join(tmdbNames(item.CreatedBy), ", ")})
	}
	details.Crew = append(details.Crew, tmdbCreatorPeople(item.CreatedBy)...)
	details.Crew = append(details.Crew, tmdbCrewPeople(item.Credits.Crew)...)
	if len(item.Networks) > 0 {
		details.Facts = append(details.Facts, Fact{Label: "Networks", Value: strings.Join(tmdbNames(item.Networks), "\n")})
	}
	details.Facts = append(details.Facts, tmdbExternalIDFacts(item.ExternalIDs)...)
	for _, cast := range item.Credits.Cast {
		name := strings.TrimSpace(cast.Name)
		if name == "" {
			continue
		}
		details.Cast = append(details.Cast, Person{
			ExternalProvider: optionalString("tmdb"),
			ExternalID:       optionalString(strconv.FormatInt(cast.ID, 10)),
			Name:             name,
			Role:             optionalString(cast.Character),
			ProfilePath:      optionalString(cast.ProfilePath),
		})
		if len(details.Cast) >= 80 {
			break
		}
	}
	return details
}

func tmdbReleaseDateFacts(info tmdbReleaseInfo) []Fact {
	selected := map[int]string{}
	for _, country := range info.Results {
		if country.Code != "US" {
			continue
		}
		for _, release := range country.ReleaseList {
			date := releaseDateOnly(release.Date)
			if date == "" {
				continue
			}
			switch release.Type {
			case 3:
				selected[3] = earliestDate(selected[3], date)
			case 4:
				selected[4] = earliestDate(selected[4], date)
			case 5:
				selected[5] = earliestDate(selected[5], date)
			}
		}
	}
	facts := []Fact{}
	if selected[3] != "" {
		facts = append(facts, Fact{Label: "Theatrical Release Date", Value: selected[3]})
	}
	if selected[4] != "" {
		facts = append(facts, Fact{Label: "Digital Release Date", Value: selected[4]})
	}
	if selected[5] != "" {
		facts = append(facts, Fact{Label: "Physical Release Date", Value: selected[5]})
	}
	return facts
}

func tmdbCertificationFacts(item tmdbDetails) []Fact {
	if value := tmdbReleaseCertification(item.ReleaseDates); value != "" {
		return []Fact{{Label: "Certification", Value: value}}
	}
	if value := tmdbContentRatingValue(item.ContentRatings); value != "" {
		return []Fact{{Label: "Certification", Value: value}}
	}
	return nil
}

func tmdbReleaseCertification(info tmdbReleaseInfo) string {
	selected := map[int]string{}
	for _, country := range info.Results {
		if country.Code != "US" {
			continue
		}
		for _, release := range country.ReleaseList {
			value := strings.TrimSpace(release.Certification)
			if value == "" {
				continue
			}
			selected[release.Type] = value
		}
	}
	for _, releaseType := range []int{3, 2, 1, 4, 5, 6} {
		if selected[releaseType] != "" {
			return selected[releaseType]
		}
	}
	return ""
}

func tmdbContentRatingValue(info tmdbContentRatings) string {
	for _, rating := range info.Results {
		if rating.Code == "US" {
			return strings.TrimSpace(rating.Rating)
		}
	}
	return ""
}

func tmdbFinancialFacts(item tmdbDetails) []Fact {
	facts := []Fact{}
	if item.Revenue > 0 {
		facts = append(facts, Fact{Label: "Revenue", Value: formatMoney(item.Revenue)})
	}
	if item.Budget > 0 {
		facts = append(facts, Fact{Label: "Budget", Value: formatMoney(item.Budget)})
	}
	return facts
}

func tmdbExternalIDFacts(ids tmdbExternalIDs) []Fact {
	facts := []Fact{}
	if value := strings.TrimSpace(ids.IMDBID); value != "" {
		facts = append(facts, Fact{Label: "IMDb ID", Value: value})
	}
	if value := strings.TrimSpace(ids.WikidataID); value != "" {
		facts = append(facts, Fact{Label: "Wikidata ID", Value: value})
	}
	if ids.TVDBID > 0 {
		facts = append(facts, Fact{Label: "TVDB ID", Value: strconv.FormatInt(ids.TVDBID, 10)})
	}
	if value := strings.TrimSpace(ids.FacebookID); value != "" {
		facts = append(facts, Fact{Label: "Facebook ID", Value: value})
	}
	if value := strings.TrimSpace(ids.InstagramID); value != "" {
		facts = append(facts, Fact{Label: "Instagram ID", Value: value})
	}
	if value := strings.TrimSpace(ids.TwitterID); value != "" {
		facts = append(facts, Fact{Label: "Twitter ID", Value: value})
	}
	return facts
}

func tmdbKeywordNames(keywords tmdbKeywords) []string {
	values := keywords.Keywords
	if len(values) == 0 {
		values = keywords.Results
	}
	return tmdbNames(values)
}
