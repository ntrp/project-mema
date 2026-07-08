package decisions

import (
	"fmt"

	"media-manager/internal/storage"
)

func releaseCriteria(item storage.MediaItem, release storage.ReleaseCandidate) ReleaseSearchCriteria {
	criteria := SearchCriteriaForQuery(item, "")
	if release.SearchKind != "" {
		criteria.Kind = release.SearchKind
	}
	criteria.SeasonID = release.SeasonID
	criteria.EpisodeID = release.EpisodeID
	criteria.SeasonNumber = release.RequestedSeason
	criteria.EpisodeNumber = release.RequestedEpisode
	return criteria
}

func parsedResourceTitle(mediaType string, parsed ParsedRelease) string {
	if mediaType == "serie" && parsed.SeriesTitle != "" {
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
		Severity:  severity,
		Details:   append([]string{}, details...),
		Parsed:    parsed,
		QualityID: parsed.QualityID,
		Quality:   parsed.Quality,
		Score:     score,
		Languages: append([]string{}, parsed.Languages...),
	}
}

func scoredReleaseMatch(
	severity string,
	parsed ParsedRelease,
	matchedMedia string,
	customScore int32,
	customContributors []ReleaseScoreContributor,
	languageScore int32,
	languageContributors []ReleaseScoreContributor,
	targetScore int32,
	targetContributors []ReleaseScoreContributor,
	details ...string,
) ReleaseMatch {
	match := releaseMatch(severity, parsed, customScore+languageScore+targetScore, details...)
	match.MatchedMedia = matchedMedia
	match.CustomFormatScore = customScore
	match.CustomFormatContributors = customContributors
	match.LanguageScore = languageScore
	match.LanguageContributors = languageContributors
	match.TargetScore = targetScore
	match.TargetContributors = targetContributors
	match.ScoreContributors = append([]ReleaseScoreContributor{}, customContributors...)
	match.ScoreContributors = append(match.ScoreContributors, languageContributors...)
	match.ScoreContributors = append(match.ScoreContributors, targetContributors...)
	return match
}

func decisionMatch(
	severity string,
	parsed ParsedRelease,
	score int32,
	matchedMedia string,
	customScore int32,
	customContributors []ReleaseScoreContributor,
	languageScore int32,
	languageContributors []ReleaseScoreContributor,
	targetScore int32,
	targetContributors []ReleaseScoreContributor,
	meta releaseMeta,
	details ...string,
) ReleaseMatch {
	match := scoredReleaseMatch(
		severity,
		parsed,
		matchedMedia,
		customScore,
		customContributors,
		languageScore,
		languageContributors,
		targetScore,
		targetContributors,
		details...,
	)
	match.RankContributors = rankContributors(parsed, score, customScore, languageScore, targetScore, meta)
	return match
}

func qualityRejection(parsed ParsedRelease, profile *storage.MediaProfile) string {
	if profile == nil {
		return ""
	}
	if parsed.QualityID == "" {
		return "Release quality could not be identified."
	}
	for _, qualityID := range profile.QualityIDs {
		if qualityID == parsed.QualityID {
			return ""
		}
	}
	return fmt.Sprintf("Quality %s is not enabled in the profile.", parsed.Quality)
}
