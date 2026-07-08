package jobs

import (
	"context"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func dedupeReleaseCandidates(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	releases []storage.ReleaseCandidateInput,
) []storage.ReleaseCandidateInput {
	byKey := map[string]storage.ReleaseCandidateInput{}
	matches := newReleaseMatchCache(item, profile, formats, languages)
	for _, release := range releases {
		release.Sources = storage.ReleaseCandidateSourcesForInput(release)
		key := releaseDedupeKey(release)
		if key == "" {
			byKey[uuid.NewString()] = release
			continue
		}
		existing, ok := byKey[key]
		if !ok {
			byKey[key] = release
			continue
		}
		if betterCandidate(matches, release, existing) {
			release.Sources = appendUniqueReleaseSources(release.Sources, existing.Sources)
			byKey[key] = release
			continue
		}
		existing.Sources = appendUniqueReleaseSources(existing.Sources, release.Sources)
		byKey[key] = existing
	}
	deduped := make([]storage.ReleaseCandidateInput, 0, len(byKey))
	for _, release := range byKey {
		deduped = append(deduped, release)
	}
	return deduped
}

func betterCandidate(
	matches releaseMatchCache,
	left storage.ReleaseCandidateInput,
	right storage.ReleaseCandidateInput,
) bool {
	leftMatch := matches.match(left)
	rightMatch := matches.match(right)
	if leftMatch.Severity != rightMatch.Severity {
		return severityRank(leftMatch.Severity) > severityRank(rightMatch.Severity)
	}
	if left.Seeders != nil && right.Seeders != nil && *left.Seeders != *right.Seeders {
		return *left.Seeders > *right.Seeders
	}
	if left.SizeBytes != right.SizeBytes {
		return left.SizeBytes > right.SizeBytes
	}
	return strings.ToLower(left.Title) < strings.ToLower(right.Title)
}

func sortReleaseCandidates(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	releases []storage.ReleaseCandidateInput,
) {
	matches := newReleaseMatchCache(item, profile, formats, languages)
	sort.SliceStable(releases, func(i, j int) bool {
		return sortedReleaseLess(matches, profile, releases[i], releases[j])
	})
	for index := range releases {
		releases[index] = releaseWithCustomFormatMatches(matches.match(releases[index]), releases[index])
	}
}

func releaseWithCustomFormatMatches(
	match decisions.ReleaseMatch,
	release storage.ReleaseCandidateInput,
) storage.ReleaseCandidateInput {
	release.CustomFormatScore = match.CustomFormatScore
	release.MatchedCustomFormats = make([]storage.ReleaseCandidateCustomFormatMatch, 0, len(match.CustomFormatContributors))
	for _, contributor := range match.CustomFormatContributors {
		release.MatchedCustomFormats = append(release.MatchedCustomFormats, storage.ReleaseCandidateCustomFormatMatch{
			Name:  contributor.Label,
			Score: contributor.Score,
		})
	}
	return release
}

func sortedReleaseLess(
	matches releaseMatchCache,
	profile *storage.MediaProfile,
	left storage.ReleaseCandidateInput,
	right storage.ReleaseCandidateInput,
) bool {
	leftMatch := matches.match(left)
	rightMatch := matches.match(right)
	if leftMatch.Severity != rightMatch.Severity {
		return severityRank(leftMatch.Severity) > severityRank(rightMatch.Severity)
	}
	if leftMatch.QualityID != rightMatch.QualityID {
		return qualityRank(leftMatch.QualityID, profile) > qualityRank(rightMatch.QualityID, profile)
	}
	if leftMatch.CustomFormatScore != rightMatch.CustomFormatScore {
		return leftMatch.CustomFormatScore > rightMatch.CustomFormatScore
	}
	if left.Seeders != nil && right.Seeders != nil && *left.Seeders != *right.Seeders {
		return *left.Seeders > *right.Seeders
	}
	if left.SizeBytes != right.SizeBytes {
		return left.SizeBytes > right.SizeBytes
	}
	return strings.ToLower(left.Title) < strings.ToLower(right.Title)
}

type releaseMatchCache struct {
	item      storage.MediaItem
	profile   *storage.MediaProfile
	formats   []storage.CustomFormat
	languages []storage.Language
	values    map[string]decisions.ReleaseMatch
}

func newReleaseMatchCache(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
) releaseMatchCache {
	return releaseMatchCache{
		item:      item,
		profile:   profile,
		formats:   formats,
		languages: languages,
		values:    map[string]decisions.ReleaseMatch{},
	}
}

func (cache releaseMatchCache) match(release storage.ReleaseCandidateInput) decisions.ReleaseMatch {
	key := releaseMatchCacheKey(release)
	if value, ok := cache.values[key]; ok {
		return value
	}
	value := decisions.EvaluateReleaseCandidateInputMatchWithLanguageContext(
		cache.item,
		release,
		cache.profile,
		cache.formats,
		cache.languages,
	)
	cache.values[key] = value
	return value
}

func releaseMatchCacheKey(release storage.ReleaseCandidateInput) string {
	values := []string{
		release.Title,
		release.DownloadURL,
		release.SearchKind,
	}
	if release.GUID != nil {
		values = append(values, *release.GUID)
	}
	if release.RequestedSeason != nil {
		values = append(values, "s", stringInt32(*release.RequestedSeason))
	}
	if release.RequestedEpisode != nil {
		values = append(values, "e", stringInt32(*release.RequestedEpisode))
	}
	return strings.Join(values, "\x00")
}

func stringInt32(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

func qualityRank(qualityID string, profile *storage.MediaProfile) int {
	if profile != nil {
		for index, value := range profile.QualityIDs {
			if value == qualityID {
				return index + 1
			}
		}
		return 0
	}
	for _, definition := range storage.QualitySizeDefinitions() {
		if definition.ID == qualityID {
			return int(definition.SortOrder)
		}
	}
	return 0
}

func releaseDecisionContext(
	ctx context.Context,
	settings *storage.SettingsStore,
	item storage.MediaItem,
) (*storage.MediaProfile, []storage.CustomFormat, []storage.Language) {
	var profile *storage.MediaProfile
	if item.QualityProfileID != nil {
		value, err := settings.GetMediaProfile(ctx, *item.QualityProfileID)
		if err != nil {
			slog.Debug("release decision profile load failed", "profileId", *item.QualityProfileID, "error", err)
		} else {
			profile = &value
		}
	}
	formats, err := settings.ListCustomFormats(ctx)
	if err != nil {
		slog.Debug("release decision custom format load failed", "error", err)
	}
	languages, err := settings.ListLanguages(ctx)
	if err != nil {
		slog.Debug("release decision language load failed", "error", err)
	}
	return profile, formats, languages
}

func releaseDedupeKey(release storage.ReleaseCandidateInput) string {
	if release.GUID != nil && strings.TrimSpace(*release.GUID) != "" {
		return "guid:" + strings.ToLower(strings.TrimSpace(*release.GUID))
	}
	if title := normalizedReleaseDedupeTitle(release.Title); title != "" && release.SizeBytes > 0 {
		return strings.Join([]string{
			"title-size",
			title,
			strconv.FormatInt(release.SizeBytes, 10),
			strings.ToLower(strings.TrimSpace(release.IndexerProtocol)),
		}, ":")
	}
	if release.InfoURL != nil && strings.TrimSpace(*release.InfoURL) != "" {
		return "info:" + strings.ToLower(strings.TrimSpace(*release.InfoURL))
	}
	if strings.TrimSpace(release.DownloadURL) != "" {
		return "download:" + strings.ToLower(strings.TrimSpace(release.DownloadURL))
	}
	return ""
}

func normalizedReleaseDedupeTitle(title string) string {
	var builder strings.Builder
	lastSeparator := false
	for _, value := range strings.ToLower(title) {
		if (value >= 'a' && value <= 'z') || (value >= '0' && value <= '9') {
			builder.WriteRune(value)
			lastSeparator = false
			continue
		}
		if !lastSeparator {
			builder.WriteByte(' ')
			lastSeparator = true
		}
	}
	return strings.Join(strings.Fields(builder.String()), " ")
}

func severityRank(severity string) int {
	switch severity {
	case "info":
		return 3
	case "warning":
		return 2
	default:
		return 1
	}
}
