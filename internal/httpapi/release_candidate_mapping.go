package httpapi

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func releaseCandidateResponse(
	item storage.MediaItem,
	release storage.ReleaseCandidate,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
) ReleaseCandidate {
	return releaseCandidateResponseWithBlock(item, release, profile, formats, languages, nil)
}

func releaseCandidateResponseWithBlock(
	item storage.MediaItem,
	release storage.ReleaseCandidate,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	block *storage.ReleaseBlocklistItem,
) ReleaseCandidate {
	return releaseCandidateResponseWithBlockAndGrabReason(item, release, profile, formats, languages, block, nil)
}

func releaseCandidateResponseWithBlockAndGrabReason(
	item storage.MediaItem,
	release storage.ReleaseCandidate,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	block *storage.ReleaseBlocklistItem,
	grabDisabledReason *string,
) ReleaseCandidate {
	var indexerID *openapi_types.UUID
	if release.IndexerID != nil {
		value := openapi_types.UUID(*release.IndexerID)
		indexerID = &value
	}
	match := decisions.EvaluateReleaseMatchWithLanguageContext(
		item,
		release,
		profile,
		formats,
		languages,
	)
	if block != nil {
		match.Severity = "error"
		match.Details = append([]string{"Release is blocklisted: " + block.Reason}, match.Details...)
	}
	return ReleaseCandidate{
		Id:                 openapi_types.UUID(release.ID),
		IndexerId:          indexerID,
		IndexerName:        release.IndexerName,
		IndexerProtocol:    IndexerProtocol(release.IndexerProtocol),
		Title:              release.Title,
		InfoUrl:            release.InfoURL,
		Guid:               release.GUID,
		SizeBytes:          release.SizeBytes,
		Seeders:            release.Seeders,
		Peers:              release.Peers,
		PublishedAt:        release.PublishedAt,
		GrabDisabledReason: grabDisabledReason,
		Sources:            releaseCandidateSourceResponses(storage.ReleaseCandidateSourcesForStored(release)),
		Match: ReleaseCandidateMatch{
			Severity:          ReleaseCandidateMatchSeverity(match.Severity),
			Details:           match.Details,
			QualityId:         match.QualityID,
			Quality:           match.Quality,
			Score:             match.Score,
			ScoreContributors: releaseScoreContributorResponses(match.ScoreContributors),
			Languages:         match.Languages,
			MatchedMedia:      match.MatchedMedia,
			MatchedSeasonId:   optionalOpenAPIUUID(match.MatchedSeasonID),
			MatchedEpisodeId:  optionalOpenAPIUUID(match.MatchedEpisodeID),
			CustomFormatScore: match.CustomFormatScore,
			CustomFormatContributors: releaseScoreContributorResponses(
				match.CustomFormatContributors,
			),
			LanguageContributors: releaseScoreContributorResponses(match.LanguageContributors),
			RankContributors:     releaseScoreContributorResponses(match.RankContributors),
			Parsed:               parsedReleaseMetadataResponse(match.Parsed, languages),
		},
	}
}

func releaseScoreContributorResponses(
	contributors []decisions.ReleaseScoreContributor,
) []ReleaseScoreContributor {
	responses := make([]ReleaseScoreContributor, 0, len(contributors))
	for _, contributor := range contributors {
		responses = append(responses, ReleaseScoreContributor{
			Label: contributor.Label,
			Score: contributor.Score,
		})
	}
	return responses
}

func releaseCandidateSourceResponses(
	sources []storage.ReleaseCandidateSource,
) []ReleaseCandidateSource {
	responses := make([]ReleaseCandidateSource, 0, len(sources))
	for _, source := range sources {
		var indexerID *openapi_types.UUID
		if source.IndexerID != nil {
			value := openapi_types.UUID(*source.IndexerID)
			indexerID = &value
		}
		responses = append(responses, ReleaseCandidateSource{
			IndexerId:       indexerID,
			IndexerName:     source.IndexerName,
			IndexerProtocol: IndexerProtocol(source.IndexerProtocol),
			Title:           source.Title,
			DownloadUrl:     source.DownloadURL,
			InfoUrl:         source.InfoURL,
			Guid:            source.GUID,
		})
	}
	return responses
}
