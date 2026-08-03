package jobs

import (
	"fmt"
	"strings"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

type automaticReleaseTarget struct {
	label      string
	candidates []storage.ReleaseCandidateInput
}

func automaticReleaseTargets(
	engine decisions.Engine,
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	releases []storage.ReleaseCandidateInput,
) ([]automaticReleaseTarget, []monitoredEpisode) {
	if item.Type != "serie" {
		if _, ok := engine.ChooseReleaseWithProfileAndLanguages(item, profile, formats, languages, releases); !ok {
			return nil, nil
		}
		return []automaticReleaseTarget{{label: item.Title, candidates: releases}}, nil
	}
	allEpisodes := monitoredAiredEpisodes(item.Seasons)
	episodes := missingMonitoredEpisodes(item)
	if len(episodes) == 0 {
		return nil, nil
	}
	seriesPacks := []storage.ReleaseCandidateInput{}
	if len(episodes) == len(allEpisodes) {
		seriesPacks = completeSeriesPackCandidates(releases)
	}
	seriesPackDecision, hasSeriesPack := engine.ChooseReleaseWithProfileAndLanguages(item, profile, formats, languages, seriesPacks)
	if hasSeriesPack && (profile == nil || profile.SeriesPackPreference != "preferEpisodes") {
		episodeDecisions := make([]decisions.ReleaseDecision, 0, len(episodes))
		for _, episode := range episodes {
			candidates := candidatesForEpisode(item, releasesForSeason(releases, episode.season), episode)
			if decision, ok := engine.ChooseReleaseWithProfileAndLanguages(item, profile, formats, languages, candidates); ok {
				episodeDecisions = append(episodeDecisions, decision)
			}
		}
		if len(episodeDecisions) < len(episodes) || packNoWorseThanEpisodes(seriesPackDecision, episodeDecisions, profile) {
			return []automaticReleaseTarget{{label: "complete series pack", candidates: seriesPacks}}, nil
		}
	}

	targets := []automaticReleaseTarget{}
	missing := []monitoredEpisode{}
	for _, seasonNumber := range episodeSeasonNumbers(episodes) {
		seasonEpisodes := episodesForSeason(episodes, seasonNumber)
		seasonReleases := releasesForSeason(releases, seasonNumber)
		packCandidates := []storage.ReleaseCandidateInput{}
		if len(seasonEpisodes) == len(episodesForSeason(allEpisodes, seasonNumber)) {
			packCandidates = seasonPackCandidates(seasonReleases)
		}
		packDecision, hasPack := engine.ChooseReleaseWithProfileAndLanguages(item, profile, formats, languages, packCandidates)

		episodeTargets := make([]automaticReleaseTarget, 0, len(seasonEpisodes))
		episodeDecisions := make([]decisions.ReleaseDecision, 0, len(seasonEpisodes))
		seasonMissing := []monitoredEpisode{}
		for _, episode := range seasonEpisodes {
			candidates := candidatesForEpisode(item, seasonReleases, episode)
			decision, ok := engine.ChooseReleaseWithProfileAndLanguages(item, profile, formats, languages, candidates)
			if !ok {
				seasonMissing = append(seasonMissing, episode)
				continue
			}
			episodeTargets = append(episodeTargets, automaticReleaseTarget{
				label:      fmt.Sprintf("S%02dE%02d", episode.season, episode.episode),
				candidates: candidates,
			})
			episodeDecisions = append(episodeDecisions, decision)
		}

		preferPack := profile == nil || profile.SeriesPackPreference != "preferEpisodes"
		if hasPack && preferPack && (len(seasonMissing) > 0 || packNoWorseThanEpisodes(packDecision, episodeDecisions, profile)) {
			targets = append(targets, automaticReleaseTarget{
				label:      fmt.Sprintf("S%02d pack", seasonNumber),
				candidates: packCandidates,
			})
			continue
		}
		targets = append(targets, episodeTargets...)
		missing = append(missing, seasonMissing...)
	}
	return targets, missing
}

func missingMonitoredEpisodes(item storage.MediaItem) []monitoredEpisode {
	existing := map[string]bool{}
	for _, file := range item.FileFacts {
		if file.EpisodeID != nil {
			existing[file.EpisodeID.String()] = true
		}
	}
	episodes := monitoredAiredEpisodes(item.Seasons)
	missing := make([]monitoredEpisode, 0, len(episodes))
	for _, episode := range episodes {
		if episode.episodeID != nil && existing[episode.episodeID.String()] {
			continue
		}
		missing = append(missing, episode)
	}
	return missing
}

func episodeSeasonNumbers(episodes []monitoredEpisode) []int32 {
	seen := map[int32]bool{}
	values := []int32{}
	for _, episode := range episodes {
		if !seen[episode.season] {
			seen[episode.season] = true
			values = append(values, episode.season)
		}
	}
	return values
}

func episodesForSeason(episodes []monitoredEpisode, season int32) []monitoredEpisode {
	values := []monitoredEpisode{}
	for _, episode := range episodes {
		if episode.season == season {
			values = append(values, episode)
		}
	}
	return values
}

func releasesForSeason(releases []storage.ReleaseCandidateInput, season int32) []storage.ReleaseCandidateInput {
	values := []storage.ReleaseCandidateInput{}
	for _, release := range releases {
		parsed := decisions.ParseReleaseFileName(release.Title)
		if release.RequestedSeason != nil && *release.RequestedSeason == season {
			values = append(values, release)
		} else if release.RequestedSeason == nil && parsed.SeasonNumber != nil && *parsed.SeasonNumber == season {
			values = append(values, release)
		}
	}
	return values
}

func completeSeriesPackCandidates(releases []storage.ReleaseCandidateInput) []storage.ReleaseCandidateInput {
	values := []storage.ReleaseCandidateInput{}
	for _, release := range releases {
		if !decisions.IsSeriesPackTitle(release.Title) {
			continue
		}
		candidate := release
		candidate.SearchKind = "serie"
		candidate.SeasonID = nil
		candidate.EpisodeID = nil
		candidate.RequestedSeason = nil
		candidate.RequestedEpisode = nil
		values = append(values, candidate)
	}
	return values
}

func seasonPackCandidates(releases []storage.ReleaseCandidateInput) []storage.ReleaseCandidateInput {
	values := []storage.ReleaseCandidateInput{}
	for _, release := range releases {
		if decisions.ParseReleaseFileName(release.Title).SeasonPack {
			values = append(values, release)
		}
	}
	return values
}

func candidatesForEpisode(item storage.MediaItem, releases []storage.ReleaseCandidateInput, episode monitoredEpisode) []storage.ReleaseCandidateInput {
	values := []storage.ReleaseCandidateInput{}
	for _, release := range releases {
		parsed := decisions.ParseReleaseFileName(release.Title)
		if parsed.SeasonPack || (parsed.SeasonNumber != nil && *parsed.SeasonNumber != episode.season) ||
			(parsed.EpisodeNumber != nil && *parsed.EpisodeNumber != episode.episode) {
			continue
		}
		if parsed.EpisodeNumber == nil && !strings.EqualFold(stringValue(item.NumberingStrategy), "anidb_absolute") {
			continue
		}
		candidate := release
		seasonNumber, episodeNumber := episode.season, episode.episode
		candidate.SearchKind = "episode"
		candidate.SeasonID = episode.seasonID
		candidate.EpisodeID = episode.episodeID
		candidate.RequestedSeason = &seasonNumber
		if item.ContentKind == "anime" && strings.EqualFold(stringValue(item.NumberingStrategy), "anidb_absolute") {
			candidate.RequestedSeason = nil
		}
		candidate.RequestedEpisode = &episodeNumber
		values = append(values, candidate)
	}
	return values
}

func packNoWorseThanEpisodes(pack decisions.ReleaseDecision, episodes []decisions.ReleaseDecision, profile *storage.MediaProfile) bool {
	if len(episodes) == 0 {
		return true
	}
	for _, episode := range episodes {
		if releaseQualityRank(pack.Match) < releaseQualityRank(episode.Match) ||
			pack.Match.CustomFormatScore < episode.Match.CustomFormatScore ||
			pack.Match.LanguageScore < episode.Match.LanguageScore || pack.Match.TargetScore < episode.Match.TargetScore ||
			decisions.ProtocolRank(pack.Release.IndexerProtocol, profile) < decisions.ProtocolRank(episode.Release.IndexerProtocol, profile) ||
			releaseAvailabilityRank(pack.Release) < releaseAvailabilityRank(episode.Release) ||
			releaseSeeders(pack.Release) < releaseSeeders(episode.Release) {
			return false
		}
	}
	return true
}

func releaseAvailabilityRank(release storage.ReleaseCandidateInput) int {
	if strings.EqualFold(release.IndexerProtocol, "torrent") && release.Seeders != nil && *release.Seeders <= 0 {
		return -1
	}
	return 0
}

func releaseSeeders(release storage.ReleaseCandidateInput) int32 {
	if release.Seeders == nil {
		return -1
	}
	return *release.Seeders
}

func releaseQualityRank(match decisions.ReleaseMatch) int32 {
	for _, contributor := range match.RankContributors {
		if strings.HasPrefix(contributor.Label, "Quality rank:") {
			return contributor.Score
		}
	}
	return 0
}

func automaticMissingEpisodeBranches(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	releases []storage.ReleaseCandidateInput,
) []releaseSearchBranch {
	_, missing := automaticReleaseTargets(decisions.NewEngine(), item, profile, formats, languages, releases)
	branches := make([]releaseSearchBranch, 0, len(missing))
	for _, episode := range missing {
		seasonNumber, episodeNumber := episode.season, episode.episode
		criteria := decisions.ReleaseSearchCriteria{
			Kind: "episode", Title: item.Title,
			SeasonID: episode.seasonID, EpisodeID: episode.episodeID,
			SeasonNumber: &seasonNumber, EpisodeNumber: &episodeNumber,
		}
		queries := decisions.SearchQueriesForCriteria(criteria, "")
		if strings.TrimSpace(episode.title) != "" {
			queries = append(queries, fmt.Sprintf("%s %s", item.Title, episode.title))
		}
		branches = append(branches, releaseSearchBranch{criteria: criteria, queries: queries})
	}
	return branches
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
