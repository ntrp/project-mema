package decisions

import (
	"regexp"
	"strings"
	"time"

	"media-manager/internal/storage"
)

type ReleaseDecision struct {
	Release storage.ReleaseCandidateInput
	Match   ReleaseMatch
}

type Engine struct {
	qualities []qualityRule
}

type qualityRule struct {
	id        string
	name      string
	sortOrder int32
	tokens    []string
}

var nonAlphaNumeric = regexp.MustCompile(`[^a-z0-9]+`)

func NewEngine() Engine {
	definitions := storage.QualitySizeDefinitions()
	qualities := make([]qualityRule, 0, len(definitions))
	for _, definition := range definitions {
		qualities = append(qualities, qualityRule{
			id:        definition.ID,
			name:      definition.Name,
			sortOrder: definition.SortOrder,
			tokens: uniqueTokens(
				normalizedToken(definition.ID),
				normalizedToken(definition.Name),
			),
		})
	}
	return Engine{qualities: qualities}
}

func (e Engine) ChooseRelease(item storage.MediaItem, candidates []storage.ReleaseCandidateInput) (ReleaseDecision, bool) {
	return e.ChooseReleaseWithProfile(item, nil, nil, candidates)
}

func (e Engine) ChooseReleaseWithProfile(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	candidates []storage.ReleaseCandidateInput,
) (ReleaseDecision, bool) {
	return e.ChooseReleaseWithProfileAndLanguages(item, profile, formats, nil, candidates)
}

func (e Engine) ChooseReleaseWithProfileAndLanguages(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	candidates []storage.ReleaseCandidateInput,
) (ReleaseDecision, bool) {
	if len(candidates) == 0 {
		return ReleaseDecision{}, false
	}

	var best storage.ReleaseCandidateInput
	var bestMatch ReleaseMatch
	found := false
	for _, candidate := range candidates {
		match := EvaluateReleaseCandidateInputMatchWithLanguageContext(
			item,
			candidate,
			profile,
			formats,
			languages,
		)
		if match.Severity == "error" {
			continue
		}
		if !found || betterRelease(candidate, match, best, bestMatch, profile) {
			best = candidate
			bestMatch = match
			found = true
		}
	}
	if !found {
		return ReleaseDecision{}, false
	}
	return ReleaseDecision{Release: best, Match: bestMatch}, true
}

func betterRelease(
	left storage.ReleaseCandidateInput,
	leftMatch ReleaseMatch,
	right storage.ReleaseCandidateInput,
	rightMatch ReleaseMatch,
	profile *storage.MediaProfile,
) bool {
	if leftMatch.QualityID != rightMatch.QualityID {
		return profileQualityScore(leftMatch.QualityID, profile) > profileQualityScore(rightMatch.QualityID, profile)
	}
	if leftMatch.CustomFormatScore != rightMatch.CustomFormatScore {
		return leftMatch.CustomFormatScore > rightMatch.CustomFormatScore
	}
	if leftProtocol := protocolRank(left.IndexerProtocol, profile); leftProtocol != protocolRank(right.IndexerProtocol, profile) {
		return leftProtocol > protocolRank(right.IndexerProtocol, profile)
	}
	if packRank := seasonPackRank(left.Title, profile) - seasonPackRank(right.Title, profile); packRank != 0 {
		return packRank > 0
	}
	leftSeeders := int32(-1)
	rightSeeders := int32(-1)
	if left.Seeders != nil {
		leftSeeders = *left.Seeders
	}
	if right.Seeders != nil {
		rightSeeders = *right.Seeders
	}
	if leftSeeders != rightSeeders {
		return leftSeeders > rightSeeders
	}
	if left.PublishedAt != nil && right.PublishedAt != nil && !left.PublishedAt.Equal(*right.PublishedAt) {
		return left.PublishedAt.After(*right.PublishedAt)
	}
	if left.SizeBytes != right.SizeBytes {
		return left.SizeBytes < right.SizeBytes
	}
	return strings.ToLower(left.Title) < strings.ToLower(right.Title)
}

func seasonPackRank(title string, profile *storage.MediaProfile) int {
	if profile == nil {
		return 0
	}
	parsed := ParseReleaseFileName(title)
	switch profile.SeriesPackPreference {
	case "preferPacks":
		if parsed.SeasonPack {
			return 1
		}
	case "preferEpisodes":
		if !parsed.SeasonPack {
			return 1
		}
	}
	return 0
}

func protocolRank(indexerProtocol string, profile *storage.MediaProfile) int {
	if profile == nil || profile.PreferredProtocol == "" || profile.PreferredProtocol == "any" {
		return 0
	}
	protocol := "torrent"
	if strings.EqualFold(indexerProtocol, "nzb") || strings.EqualFold(indexerProtocol, "usenet") {
		protocol = "usenet"
	}
	if protocol == profile.PreferredProtocol {
		return 1
	}
	return 0
}

func (e Engine) detectQuality(title string) qualityRule {
	normalizedTitle := normalizedToken(title)
	best := qualityRule{}
	for _, quality := range e.qualities {
		for _, token := range quality.tokens {
			if token == "" || !strings.Contains(normalizedTitle, token) {
				continue
			}
			if quality.sortOrder > best.sortOrder {
				best = quality
			}
		}
	}
	return best
}

func EvaluateReleaseCandidateInputMatchWithContext(
	item storage.MediaItem,
	release storage.ReleaseCandidateInput,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
) ReleaseMatch {
	return EvaluateReleaseCandidateInputMatchWithLanguageContext(item, release, profile, formats, nil)
}

func EvaluateReleaseCandidateInputMatchWithLanguageContext(
	item storage.MediaItem,
	release storage.ReleaseCandidateInput,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
) ReleaseMatch {
	return EvaluateReleaseMatchWithLanguageContext(
		item,
		releaseCandidateFromInput(release),
		profile,
		formats,
		languages,
	)
}

func releaseCandidateFromInput(release storage.ReleaseCandidateInput) storage.ReleaseCandidate {
	return storage.ReleaseCandidate{
		MediaItemID:      release.MediaItemID,
		SeasonID:         release.SeasonID,
		EpisodeID:        release.EpisodeID,
		IndexerID:        release.IndexerID,
		IndexerName:      release.IndexerName,
		IndexerProtocol:  release.IndexerProtocol,
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
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	}
}

func normalizedToken(value string) string {
	return nonAlphaNumeric.ReplaceAllString(strings.ToLower(value), "")
}

func uniqueTokens(values ...string) []string {
	seen := map[string]struct{}{}
	tokens := []string{}
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		tokens = append(tokens, value)
	}
	return tokens
}
