package jobs

import (
	"encoding/json"
	"strings"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

type AudioFulfillmentNeed struct {
	LanguageID           string
	TargetCodec          string
	TargetChannels       []string
	MinimumBitrateKbps   int32
	PreferredBitrateKbps int32
	Query                string
}

type AudioFulfillmentCandidate struct {
	Need         AudioFulfillmentNeed
	ReleaseTitle string
	IndexerName  string
	DownloadURL  string
	Score        int32
	Metadata     string
}

func PlanAudioFulfillment(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	releases []storage.ReleaseCandidateInput,
	languages []storage.Language,
) []AudioFulfillmentCandidate {
	if profile == nil {
		return nil
	}
	needs := AudioFulfillmentNeeds(item, *profile)
	candidates := []AudioFulfillmentCandidate{}
	for _, need := range needs {
		if candidate, ok := bestAudioFulfillmentCandidate(item, profile, releases, languages, need); ok {
			candidates = append(candidates, candidate)
		}
	}
	return candidates
}

func AudioFulfillmentNeeds(item storage.MediaItem, profile storage.MediaProfile) []AudioFulfillmentNeed {
	needs := []AudioFulfillmentNeed{}
	for _, target := range profile.AudioTargets {
		if audioTargetSatisfied(item, target) {
			continue
		}
		need := AudioFulfillmentNeed{
			LanguageID:           target.LanguageID,
			TargetChannels:       append([]string{}, target.TargetChannels...),
			MinimumBitrateKbps:   int32PtrValue(target.MinimumBitrateKbps),
			PreferredBitrateKbps: int32PtrValue(target.PreferredBitrateKbps),
		}
		if target.TargetCodec != nil {
			need.TargetCodec = strings.TrimSpace(*target.TargetCodec)
		}
		need.Query = audioFulfillmentQuery(item, need)
		needs = append(needs, need)
	}
	return needs
}

func bestAudioFulfillmentCandidate(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	releases []storage.ReleaseCandidateInput,
	languages []storage.Language,
	need AudioFulfillmentNeed,
) (AudioFulfillmentCandidate, bool) {
	var best AudioFulfillmentCandidate
	for _, release := range releases {
		match := decisions.EvaluateReleaseMatchWithLanguageContext(item, releaseCandidate(release), profile, nil, languages)
		if match.Severity == "error" || match.TargetScore == 0 {
			continue
		}
		candidate := AudioFulfillmentCandidate{
			Need:         need,
			ReleaseTitle: release.Title,
			IndexerName:  release.IndexerName,
			DownloadURL:  release.DownloadURL,
			Score:        match.Score,
			Metadata:     audioFulfillmentMetadata(release, need),
		}
		if candidate.Score > best.Score {
			best = candidate
		}
	}
	return best, best.ReleaseTitle != ""
}

func AudioFulfillmentSourceInput(
	candidate AudioFulfillmentCandidate,
	sourcePath string,
	streamInventory string,
) storage.MediaComponentSourceInput {
	releaseTitle := candidate.ReleaseTitle
	releaseName := candidate.ReleaseTitle
	releaseID := candidate.DownloadURL
	metadata := candidate.Metadata
	return storage.MediaComponentSourceInput{
		SourceRole:      "audio",
		SourceFilePath:  sourcePath,
		ReleaseTitle:    &releaseTitle,
		ReleaseName:     &releaseName,
		ReleaseID:       &releaseID,
		SourceMetadata:  &metadata,
		StreamInventory: streamInventory,
	}
}

func audioFulfillmentQuery(item storage.MediaItem, need AudioFulfillmentNeed) string {
	parts := []string{decisions.SearchQueryForMediaItem(item), need.LanguageID}
	if need.TargetCodec != "" {
		parts = append(parts, need.TargetCodec)
	}
	if len(need.TargetChannels) > 0 {
		parts = append(parts, need.TargetChannels[0])
	}
	return strings.Join(nonEmptyStrings(parts), " ")
}

func audioTargetSatisfied(item storage.MediaItem, target storage.MediaProfileAudioTarget) bool {
	for _, source := range item.ComponentSources {
		for _, artifact := range source.Artifacts {
			if artifact.StreamType == "audio" && languageMatches(artifact.Language, target.LanguageID) {
				return true
			}
		}
		for _, stream := range audioInventoryStreams(source.StreamInventory) {
			if stream.Type == "audio" && audioStreamMatchesTarget(stream, target) {
				return true
			}
		}
	}
	return false
}

type audioInventoryStream struct {
	Type        string `json:"type"`
	Language    string `json:"language"`
	Codec       string `json:"codec"`
	Channels    string `json:"channels"`
	BitrateKbps int32  `json:"bitrateKbps"`
}

func audioInventoryStreams(payload string) []audioInventoryStream {
	var list []audioInventoryStream
	if err := json.Unmarshal([]byte(payload), &list); err == nil {
		return list
	}
	var wrapped struct {
		Streams []audioInventoryStream `json:"streams"`
	}
	if err := json.Unmarshal([]byte(payload), &wrapped); err != nil {
		return nil
	}
	return wrapped.Streams
}

func audioStreamMatchesTarget(stream audioInventoryStream, target storage.MediaProfileAudioTarget) bool {
	if !strings.EqualFold(strings.TrimSpace(stream.Language), strings.TrimSpace(target.LanguageID)) {
		return false
	}
	if target.TargetCodec != nil && normalizeJobAudioCodec(stream.Codec) != normalizeJobAudioCodec(*target.TargetCodec) {
		return false
	}
	if len(target.TargetChannels) > 0 && !stringListHas(target.TargetChannels, stream.Channels) {
		return false
	}
	if target.MinimumBitrateKbps != nil && stream.BitrateKbps > 0 && stream.BitrateKbps < *target.MinimumBitrateKbps {
		return false
	}
	return true
}

func audioFulfillmentMetadata(release storage.ReleaseCandidateInput, need AudioFulfillmentNeed) string {
	payload := map[string]any{
		"kind":        "audioFulfillment",
		"release":     release.Title,
		"indexer":     release.IndexerName,
		"language":    need.LanguageID,
		"targetCodec": need.TargetCodec,
	}
	data, _ := json.Marshal(payload)
	return string(data)
}

func releaseCandidate(input storage.ReleaseCandidateInput) storage.ReleaseCandidate {
	return storage.ReleaseCandidate{Title: input.Title, IndexerName: input.IndexerName, DownloadURL: input.DownloadURL}
}

func languageMatches(value *string, target string) bool {
	return value != nil && strings.EqualFold(strings.TrimSpace(*value), strings.TrimSpace(target))
}

func int32PtrValue(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
}

func stringListHas(values []string, candidate string) bool {
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(candidate)) {
			return true
		}
	}
	return false
}

func nonEmptyStrings(values []string) []string {
	result := []string{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}
