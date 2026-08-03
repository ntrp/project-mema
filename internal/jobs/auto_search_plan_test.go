package jobs

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func TestAutomaticReleaseTargetsReturnNoMovieTargetWithoutAcceptableRelease(t *testing.T) {
	targets, _ := automaticReleaseTargets(
		decisions.NewEngine(), storage.MediaItem{Type: "movie", Title: "Movie"}, nil, nil, nil, nil,
	)
	if len(targets) != 0 {
		t.Fatalf("targets=%#v", targets)
	}
}

func TestAutomaticReleaseTargetsChooseEveryEpisodeWhenTheyOutrankPack(t *testing.T) {
	item := automaticPlanSeries(2)
	season := int32(1)
	profile := &storage.MediaProfile{QualityIDs: []string{"webdl-720p", "webdl-1080p"}}
	releases := []storage.ReleaseCandidateInput{
		{Title: "The.Show.S01.720p.WEBDL", SearchKind: "season", RequestedSeason: &season},
		{Title: "The.Show.S01E01.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season},
		{Title: "The.Show.S01E02.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season},
	}

	targets, missing := automaticReleaseTargets(decisions.NewEngine(), item, profile, nil, nil, releases)

	if len(missing) != 0 || len(targets) != 2 || targets[0].label != "S01E01" || targets[1].label != "S01E02" {
		t.Fatalf("targets=%#v missing=%#v", targets, missing)
	}
}

func TestAutomaticReleaseTargetsChoosePackWhenItDoesNotLoseOnQualityOrScore(t *testing.T) {
	item := automaticPlanSeries(2)
	season := int32(1)
	profile := &storage.MediaProfile{QualityIDs: []string{"webdl-720p", "webdl-1080p"}}
	releases := []storage.ReleaseCandidateInput{
		{Title: "The.Show.S01.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season},
		{Title: "The.Show.S01E01.720p.WEBDL", SearchKind: "season", RequestedSeason: &season},
		{Title: "The.Show.S01E02.720p.WEBDL", SearchKind: "season", RequestedSeason: &season},
	}

	targets, missing := automaticReleaseTargets(decisions.NewEngine(), item, profile, nil, nil, releases)

	if len(missing) != 0 || len(targets) != 1 || targets[0].label != "S01 pack" {
		t.Fatalf("targets=%#v missing=%#v", targets, missing)
	}
}

func TestAutomaticReleaseTargetsPreferEqualQualityCompleteSeriesPack(t *testing.T) {
	item := automaticPlanSeries(2)
	item.Title = "Death Note"
	season := int32(1)
	profile := &storage.MediaProfile{QualityIDs: []string{"webdl-1080p", "bluray-1080p"}}
	releases := []storage.ReleaseCandidateInput{
		{Title: "Death Note The Complete TV Series (2006&ndash;2007) 1080P BluRay", SearchKind: "season", RequestedSeason: &season},
		{Title: "Death.Note.S01E01.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season},
		{Title: "Death.Note.S01E02.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season},
	}

	targets, missing := automaticReleaseTargets(decisions.NewEngine(), item, profile, nil, nil, releases)

	if len(missing) != 0 || len(targets) != 1 || targets[0].label != "complete series pack" {
		t.Fatalf("targets=%#v missing=%#v", targets, missing)
	}
}

func TestAutomaticReleaseTargetsRejectUnseededPackAgainstAvailableEpisode(t *testing.T) {
	item := automaticPlanSeries(1)
	season, zeroSeeders := int32(1), int32(0)
	profile := &storage.MediaProfile{QualityIDs: []string{"webdl-1080p"}, PreferredProtocol: "usenet"}
	releases := []storage.ReleaseCandidateInput{
		{Title: "The.Show.S01.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season, IndexerProtocol: "torrent", Seeders: &zeroSeeders},
		{Title: "The.Show.S01E01.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season, IndexerProtocol: "usenet"},
	}

	targets, _ := automaticReleaseTargets(decisions.NewEngine(), item, profile, nil, nil, releases)

	if len(targets) != 1 || targets[0].label != "S01E01" {
		t.Fatalf("targets=%#v", targets)
	}
}

func TestAutomaticReleaseTargetsDoNotRedownloadExistingEpisodeThroughPack(t *testing.T) {
	item := automaticPlanSeries(2)
	item.FileFacts = []storage.MediaFileFact{{EpisodeID: item.Seasons[0].Episodes[0].ID}}
	season := int32(1)
	releases := []storage.ReleaseCandidateInput{
		{Title: "The.Show.S01.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season},
		{Title: "The.Show.S01E02.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season},
	}

	targets, missing := automaticReleaseTargets(decisions.NewEngine(), item, nil, nil, nil, releases)

	if len(missing) != 0 || len(targets) != 1 || targets[0].label != "S01E02" {
		t.Fatalf("targets=%#v missing=%#v", targets, missing)
	}
}

func TestAutomaticReleaseTargetsMatchAniDBAbsoluteEpisode(t *testing.T) {
	item := automaticPlanSeries(2)
	item.ContentKind = "anime"
	strategy := "anidb_absolute"
	item.NumberingStrategy = &strategy
	season := int32(1)
	releases := []storage.ReleaseCandidateInput{{
		Title: "The Show - 2 1080p WEBDL", SearchKind: "season", RequestedSeason: &season,
	}}

	targets, missing := automaticReleaseTargets(decisions.NewEngine(), item, nil, nil, nil, releases)

	if len(targets) != 1 || targets[0].label != "S01E02" || len(missing) != 1 || missing[0].episode != 1 {
		t.Fatalf("targets=%#v missing=%#v", targets, missing)
	}
}

func TestAutomaticReleaseTargetsReportEpisodesWithoutPackOrCandidate(t *testing.T) {
	item := automaticPlanSeries(2)
	season := int32(1)
	releases := []storage.ReleaseCandidateInput{{
		Title: "The.Show.S01E01.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season,
	}}

	targets, missing := automaticReleaseTargets(decisions.NewEngine(), item, nil, nil, nil, releases)

	if len(targets) != 1 || targets[0].label != "S01E01" || len(missing) != 1 || missing[0].episode != 2 {
		t.Fatalf("targets=%#v missing=%#v", targets, missing)
	}
}

func TestAutomaticMissingEpisodeBranchesSearchOnlyUncoveredEpisodes(t *testing.T) {
	item := automaticPlanSeries(2)
	season := int32(1)
	releases := []storage.ReleaseCandidateInput{{
		Title: "The.Show.S01E01.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season,
	}}

	branches := automaticMissingEpisodeBranches(item, nil, nil, nil, releases)

	if len(branches) != 1 || branches[0].criteria.EpisodeNumber == nil || *branches[0].criteria.EpisodeNumber != 2 {
		t.Fatalf("branches=%#v", branches)
	}
}

func TestAutomaticMissingEpisodeBranchesSkipIndividualSearchWhenPackCoversSeason(t *testing.T) {
	item := automaticPlanSeries(2)
	season := int32(1)
	releases := []storage.ReleaseCandidateInput{{
		Title: "The.Show.S01.1080p.WEBDL", SearchKind: "season", RequestedSeason: &season,
	}}

	if branches := automaticMissingEpisodeBranches(item, nil, nil, nil, releases); len(branches) != 0 {
		t.Fatalf("branches=%#v", branches)
	}
}

func automaticPlanSeries(episodeCount int) storage.MediaItem {
	seasonID := uuid.New()
	episodes := make([]storage.MediaEpisode, 0, episodeCount)
	for number := 1; number <= episodeCount; number++ {
		episodeID := uuid.New()
		episodes = append(episodes, storage.MediaEpisode{
			ID: &episodeID, Name: "Episode", EpisodeNumber: int32(number), Monitored: true,
		})
	}
	return storage.MediaItem{
		Type: "serie", Title: "The Show",
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{Seasons: []storage.MediaSeason{{
			ID: &seasonID, Name: "Season 1", SeasonNumber: 1, Monitored: true, Episodes: episodes,
		}}},
	}
}
