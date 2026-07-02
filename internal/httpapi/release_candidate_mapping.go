package httpapi

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func releaseCandidateResponse(item storage.MediaItem, release storage.ReleaseCandidate) ReleaseCandidate {
	var indexerID *openapi_types.UUID
	if release.IndexerID != nil {
		value := openapi_types.UUID(*release.IndexerID)
		indexerID = &value
	}
	match := decisions.EvaluateReleaseMatch(item, release)
	return ReleaseCandidate{
		Id:          openapi_types.UUID(release.ID),
		IndexerId:   indexerID,
		IndexerName: release.IndexerName,
		IndexerType: IndexerType(release.IndexerType),
		Title:       release.Title,
		InfoUrl:     release.InfoURL,
		Guid:        release.GUID,
		SizeBytes:   release.SizeBytes,
		Seeders:     release.Seeders,
		Peers:       release.Peers,
		PublishedAt: release.PublishedAt,
		Match: ReleaseCandidateMatch{
			Severity:          ReleaseCandidateMatchSeverity(match.Severity),
			Details:           match.Details,
			QualityId:         match.QualityID,
			Quality:           match.Quality,
			Score:             match.Score,
			ScoreContributors: releaseScoreContributorResponses(match.ScoreContributors),
			Languages:         match.Languages,
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
