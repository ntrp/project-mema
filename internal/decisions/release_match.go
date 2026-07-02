package decisions

import (
	"fmt"
	"strconv"
	"strings"

	"media-manager/internal/storage"
)

type ReleaseMatch struct {
	Severity                 string
	Details                  []string
	Parsed                   ParsedRelease
	QualityID                string
	Quality                  string
	Score                    int32
	ScoreContributors        []ReleaseScoreContributor
	Languages                []string
	MatchedMedia             string
	CustomFormatScore        int32
	CustomFormatContributors []ReleaseScoreContributor
	LanguageContributors     []ReleaseScoreContributor
	RankContributors         []ReleaseScoreContributor
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
	return EvaluateReleaseMatchWithContext(item, release, nil, nil)
}

func EvaluateReleaseMatchWithContext(
	item storage.MediaItem,
	release storage.ReleaseCandidate,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
) ReleaseMatch {
	return EvaluateReleaseMatchWithLanguageContext(item, release, profile, formats, nil)
}

func EvaluateReleaseMatchWithLanguageContext(
	item storage.MediaItem,
	release storage.ReleaseCandidate,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
) ReleaseMatch {
	criteria := releaseCriteria(item, release)
	parsed := ParseReleaseFileName(release.Title)
	context := ReleaseEvaluationContext{
		Item:      item,
		Profile:   profile,
		Formats:   formats,
		Languages: languages,
	}
	return evaluateParsedRelease(context, criteria, parsed, releaseMeta{
		Title:       release.Title,
		IndexerType: release.IndexerType,
		SizeBytes:   release.SizeBytes,
		Seeders:     release.Seeders,
		Peers:       release.Peers,
		PublishedAt: release.PublishedAt,
	})
}

type ReleaseEvaluationContext struct {
	Item      storage.MediaItem
	Profile   *storage.MediaProfile
	Formats   []storage.CustomFormat
	Languages []storage.Language
}

type releaseMeta struct {
	Title       string
	IndexerType string
	SizeBytes   int64
	Seeders     *int32
	Peers       *int32
	PublishedAt any
}

func evaluateParsedRelease(
	context ReleaseEvaluationContext,
	criteria ReleaseSearchCriteria,
	parsed ParsedRelease,
	meta releaseMeta,
) ReleaseMatch {
	item := context.Item
	parsed = applyLanguageCatalog(parsed, context.Languages)
	score := profileQualityScore(parsed.QualityID, context.Profile)
	details := []string{}
	matchedMedia := parsedResourceTitle(item.Type, parsed)
	customScore, customContributors := customFormatScore(parsed, context.Profile, context.Formats)
	languageScore, languageContributors, languageReject := languageScore(
		parsed,
		context.Profile,
		context.Languages,
	)

	if !resourceTitleMatches(criteria.Title, matchedMedia, meta.Title) {
		return scoredReleaseMatch("error", parsed, matchedMedia, customScore, customContributors, languageScore, languageContributors, "Does not match this series/movie.")
	}
	if yearMismatch(criteria.Year, parsed.Year) {
		return scoredReleaseMatch("error", parsed, matchedMedia, customScore, customContributors, languageScore, languageContributors, "Does not match this series/movie.", fmt.Sprintf("Release year is %s.", parsed.Year))
	}
	if reason := criteriaMismatch(criteria, parsed); reason != "" {
		return scoredReleaseMatch("error", parsed, matchedMedia, customScore, customContributors, languageScore, languageContributors, "Does not match this series/movie.", reason)
	}
	if reason := qualityRejection(parsed, context.Profile); reason != "" {
		return scoredReleaseMatch("error", parsed, matchedMedia, customScore, customContributors, languageScore, languageContributors, reason)
	}

	if languageReject != "" {
		return scoredReleaseMatch("error", parsed, matchedMedia, customScore, customContributors, languageScore, languageContributors, languageReject)
	}
	if context.Profile != nil && customScore < context.Profile.MinimumCustomFormatScore {
		return scoredReleaseMatch("error", parsed, matchedMedia, customScore, customContributors, languageScore, languageContributors, "Custom format score is below the profile minimum.")
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
		return decisionMatch("warning", parsed, score, matchedMedia, customScore, customContributors, languageScore, languageContributors, meta, details...)
	}
	if score <= currentScore && currentScore > 0 {
		return decisionMatch("warning", parsed, score, matchedMedia, customScore, customContributors, languageScore, languageContributors, meta, details...)
	}
	return decisionMatch("info", parsed, score, matchedMedia, customScore, customContributors, languageScore, languageContributors, meta, details...)
}

func EvaluateReleaseCandidateInputMatch(
	item storage.MediaItem,
	release storage.ReleaseCandidateInput,
) ReleaseMatch {
	return EvaluateReleaseCandidateInputMatchWithContext(item, release, nil, nil)
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
