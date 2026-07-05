package jobs

import (
	"context"
	"strings"

	"media-manager/internal/decisions"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

func unseenRSSReleases(config storage.Indexer, releases []indexers.Release) ([]indexers.Release, bool) {
	if config.RSSMarkerGUID == nil && config.RSSMarkerDownloadURL == nil && config.RSSMarkerPublishedAt == nil {
		return releases, true
	}
	unseen := []indexers.Release{}
	for _, release := range releases {
		if rssMarkerMatches(config, release) {
			return unseen, true
		}
		unseen = append(unseen, release)
	}
	return unseen, len(releases) == 0
}

func rssMarkerMatches(config storage.Indexer, release indexers.Release) bool {
	if config.RSSMarkerGUID != nil && markerText(*config.RSSMarkerGUID) != "" &&
		markerText(*config.RSSMarkerGUID) == markerText(release.GUID) {
		return true
	}
	if config.RSSMarkerDownloadURL != nil && markerText(*config.RSSMarkerDownloadURL) != "" &&
		markerText(*config.RSSMarkerDownloadURL) == markerText(release.DownloadURL) {
		return true
	}
	if config.RSSMarkerPublishedAt != nil && release.PublishedAt != nil &&
		release.PublishedAt.Equal(*config.RSSMarkerPublishedAt) {
		return true
	}
	return false
}

func newestRSSMarker(releases []indexers.Release) storage.RSSMarkerInput {
	if len(releases) == 0 {
		return storage.RSSMarkerInput{}
	}
	newest := releases[0]
	for _, release := range releases[1:] {
		if release.PublishedAt == nil {
			continue
		}
		if newest.PublishedAt == nil || release.PublishedAt.After(*newest.PublishedAt) {
			newest = release
		}
	}
	return storage.RSSMarkerInput{
		PublishedAt: newest.PublishedAt,
		GUID:        optionalString(newest.GUID),
		DownloadURL: optionalString(newest.DownloadURL),
	}
}

func rssReleaseCandidate(
	ctx context.Context,
	settings *storage.SettingsStore,
	item storage.MediaItem,
	release indexers.Release,
) (storage.ReleaseCandidateInput, bool) {
	candidate := releaseCandidateInput(item.ID, release, decisions.ReleaseSearchCriteria{Kind: "rss"})
	parsed := decisions.ParseReleaseFileName(release.Title)
	candidate.RequestedSeason = parsed.SeasonNumber
	candidate.RequestedEpisode = parsed.EpisodeNumber
	profile, formats, languages := releaseDecisionContext(ctx, settings, item)
	check := candidate
	check.SearchKind = rssEvaluationKind(item, parsed)
	match := decisions.EvaluateReleaseCandidateInputMatchWithLanguageContext(item, check, profile, formats, languages)
	if match.Severity == "error" {
		return storage.ReleaseCandidateInput{}, false
	}
	return candidate, true
}

func rssEvaluationKind(item storage.MediaItem, parsed decisions.ParsedRelease) string {
	if item.Type != "series" {
		return "movie"
	}
	if parsed.SeasonNumber != nil && parsed.EpisodeNumber != nil {
		return "episode"
	}
	if parsed.SeasonNumber != nil {
		return "season"
	}
	return "series"
}

func markerText(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
