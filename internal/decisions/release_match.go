package decisions

import (
	"fmt"
	"strconv"
	"strings"

	"media-manager/internal/storage"
)

type ReleaseMatch struct {
	Severity          string
	Details           []string
	QualityID         string
	Quality           string
	Score             int32
	ScoreContributors []ReleaseScoreContributor
	Languages         []string
}

type ReleaseScoreContributor struct {
	Label string
	Score int32
}

type ReleaseSearchCriteria struct {
	Kind          string
	Title         string
	Year          *int32
	SeasonNumber  *int32
	EpisodeNumber *int32
}

func SearchCriteriaForQuery(item storage.MediaItem, query string) ReleaseSearchCriteria {
	criteria := ReleaseSearchCriteria{Title: item.Title, Year: item.Year}
	if item.Type == "movie" {
		criteria.Kind = "movie"
		return criteria
	}
	season, episode := detectSeasonEpisode(query)
	criteria.SeasonNumber = season
	criteria.EpisodeNumber = episode
	switch {
	case season != nil && episode != nil:
		criteria.Kind = "episode"
	case season != nil:
		criteria.Kind = "season"
	default:
		criteria.Kind = "series"
	}
	return criteria
}

func SearchQueryForMediaItem(item storage.MediaItem) string {
	if item.Year == nil {
		return item.Title
	}
	return fmt.Sprintf("%s %d", item.Title, *item.Year)
}

func SearchQueriesForCriteria(criteria ReleaseSearchCriteria, original string) []string {
	queries := []string{}
	addQuery := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		for _, existing := range queries {
			if strings.EqualFold(existing, value) {
				return
			}
		}
		queries = append(queries, value)
	}
	addQuery(original)
	switch criteria.Kind {
	case "season":
		if criteria.SeasonNumber != nil {
			addQuery(fmt.Sprintf("%s s%d", criteria.Title, *criteria.SeasonNumber))
			addQuery(fmt.Sprintf("%s S%s", criteria.Title, padded(*criteria.SeasonNumber, 2)))
		}
	case "episode":
		if criteria.SeasonNumber != nil && criteria.EpisodeNumber != nil {
			addQuery(fmt.Sprintf("%s s%de%d", criteria.Title, *criteria.SeasonNumber, *criteria.EpisodeNumber))
			addQuery(fmt.Sprintf("%s S%sE%s", criteria.Title, padded(*criteria.SeasonNumber, 2), padded(*criteria.EpisodeNumber, 2)))
		}
	}
	return queries
}

func EvaluateReleaseMatch(item storage.MediaItem, release storage.ReleaseCandidate) ReleaseMatch {
	criteria := releaseCriteria(item, release)
	parsed := ParseReleaseFileName(release.Title)
	score := releaseQualityScore(parsed.QualityID)
	details := []string{}

	if !resourceTitleMatches(criteria.Title, parsedResourceTitle(item.Type, parsed), release.Title) {
		return releaseMatch("error", parsed, score, "Does not match this series/movie.")
	}
	if yearMismatch(criteria.Year, parsed.Year) {
		return releaseMatch("error", parsed, score, "Does not match this series/movie.", fmt.Sprintf("Release year is %s.", parsed.Year))
	}
	if reason := criteriaMismatch(criteria, parsed); reason != "" {
		return releaseMatch("error", parsed, score, "Does not match this series/movie.", reason)
	}

	details = append(details, "Matches the requested resource.")
	currentScore := currentQualityScore(item)
	switch {
	case score > currentScore:
		details = append(details, "Score is higher than the current file.")
	case currentScore > 0:
		details = append(details, "Score is lower than or equal to the current file.")
	default:
		details = append(details, "No current file score is available.")
	}
	if criteria.Kind == "episode" && parsed.SeasonPack {
		details = append(details, "This is a whole season release, but the search requested an episode.")
		return releaseMatch("warning", parsed, score, details...)
	}
	if score <= currentScore && currentScore > 0 {
		return releaseMatch("warning", parsed, score, details...)
	}
	return releaseMatch("info", parsed, score, details...)
}

func EvaluateReleaseCandidateInputMatch(
	item storage.MediaItem,
	release storage.ReleaseCandidateInput,
) ReleaseMatch {
	return EvaluateReleaseMatch(item, storage.ReleaseCandidate{
		MediaItemID:      release.MediaItemID,
		IndexerID:        release.IndexerID,
		IndexerName:      release.IndexerName,
		IndexerType:      release.IndexerType,
		Title:            release.Title,
		DownloadURL:      release.DownloadURL,
		InfoURL:          release.InfoURL,
		GUID:             release.GUID,
		SizeBytes:        release.SizeBytes,
		Seeders:          release.Seeders,
		Peers:            release.Peers,
		PublishedAt:      release.PublishedAt,
		SearchKind:       release.SearchKind,
		RequestedSeason:  release.RequestedSeason,
		RequestedEpisode: release.RequestedEpisode,
	})
}

func releaseCriteria(item storage.MediaItem, release storage.ReleaseCandidate) ReleaseSearchCriteria {
	criteria := SearchCriteriaForQuery(item, "")
	if release.SearchKind != "" {
		criteria.Kind = release.SearchKind
	}
	criteria.SeasonNumber = release.RequestedSeason
	criteria.EpisodeNumber = release.RequestedEpisode
	return criteria
}

func parsedResourceTitle(mediaType string, parsed ParsedRelease) string {
	if mediaType == "series" && parsed.SeriesTitle != "" {
		return parsed.SeriesTitle
	}
	return parsed.MovieTitle
}

func criteriaMismatch(criteria ReleaseSearchCriteria, parsed ParsedRelease) string {
	switch criteria.Kind {
	case "season":
		if criteria.SeasonNumber != nil && !sameInt32(criteria.SeasonNumber, parsed.SeasonNumber) {
			return fmt.Sprintf("Release season does not match S%s.", padded(*criteria.SeasonNumber, 2))
		}
	case "episode":
		if criteria.SeasonNumber != nil && !sameInt32(criteria.SeasonNumber, parsed.SeasonNumber) {
			return fmt.Sprintf("Release season does not match S%s.", padded(*criteria.SeasonNumber, 2))
		}
		if parsed.SeasonPack {
			return ""
		}
		if criteria.EpisodeNumber != nil && !sameInt32(criteria.EpisodeNumber, parsed.EpisodeNumber) {
			return fmt.Sprintf("Release episode does not match E%s.", padded(*criteria.EpisodeNumber, 2))
		}
	}
	return ""
}

func releaseMatch(severity string, parsed ParsedRelease, score int32, details ...string) ReleaseMatch {
	return ReleaseMatch{
		Severity:          severity,
		Details:           append([]string{}, details...),
		QualityID:         parsed.QualityID,
		Quality:           parsed.Quality,
		Score:             score,
		ScoreContributors: releaseScoreContributors(parsed, score),
		Languages:         append([]string{}, parsed.Languages...),
	}
}

func releaseScoreContributors(parsed ParsedRelease, score int32) []ReleaseScoreContributor {
	label := "Quality"
	if parsed.Quality != "" {
		label = fmt.Sprintf("Quality: %s", parsed.Quality)
	}
	return []ReleaseScoreContributor{{Label: label, Score: score}}
}

func currentQualityScore(item storage.MediaItem) int32 {
	var best int32
	for _, path := range item.FilePaths {
		if score := releaseQualityScore(ParseReleaseFileName(path).QualityID); score > best {
			best = score
		}
	}
	return best
}

func releaseQualityScore(qualityID string) int32 {
	for _, definition := range storage.QualitySizeDefinitions() {
		if definition.ID == qualityID {
			return definition.SortOrder * 100
		}
	}
	return 0
}

func yearMismatch(expected *int32, actual string) bool {
	if expected == nil || strings.TrimSpace(actual) == "" {
		return false
	}
	return strconv.Itoa(int(*expected)) != actual
}

func resourceTitleMatches(expected string, parsedTitle string, releaseTitle string) bool {
	expectedTitle := normalizedResourceTitle(expected)
	if expectedTitle == "" {
		return false
	}
	candidateTitle := normalizedResourceTitle(parsedTitle)
	if candidateTitle == "" {
		candidateTitle = normalizedResourceTitle(releaseTitle)
	}
	return expectedTitle == candidateTitle
}

func normalizedResourceTitle(title string) string {
	return normalizedToken(cleanReleaseResourceTitle(title))
}

func sameInt32(expected *int32, actual *int32) bool {
	return expected != nil && actual != nil && *expected == *actual
}

func padded(value int32, width int) string {
	return fmt.Sprintf("%0*d", width, value)
}
