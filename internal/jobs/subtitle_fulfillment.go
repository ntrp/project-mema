package jobs

import (
	"encoding/json"
	"strings"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

type SubtitleFulfillmentNeed struct {
	LanguageID string
	Source     string
	Formats    []string
	Mode       string
	Query      string
}

type SubtitleFulfillmentCandidate struct {
	Need         SubtitleFulfillmentNeed
	ReleaseTitle string
	IndexerName  string
	DownloadURL  string
	Metadata     string
}

func SubtitleFulfillmentNeeds(item storage.MediaItem) []SubtitleFulfillmentNeed {
	needs := []SubtitleFulfillmentNeed{}
	for _, target := range item.SubtitleTargets {
		for _, mode := range subtitleFulfillmentModes(item, target) {
			need := SubtitleFulfillmentNeed{
				LanguageID: target.LanguageID,
				Source:     target.Source,
				Formats:    normalizedSubtitleFormats(target.Formats),
				Mode:       mode,
			}
			need.Query = subtitleFulfillmentQuery(item, need)
			needs = append(needs, need)
		}
	}
	return needs
}

func PlanSubtitleFulfillment(
	item storage.MediaItem,
	releases []storage.ReleaseCandidateInput,
) []SubtitleFulfillmentCandidate {
	planned := []SubtitleFulfillmentCandidate{}
	for _, need := range SubtitleFulfillmentNeeds(item) {
		if need.Mode != "alternateRelease" {
			continue
		}
		if candidate, ok := bestSubtitleFulfillmentCandidate(item, releases, need); ok {
			planned = append(planned, candidate)
		}
	}
	return planned
}

func SubtitleFulfillmentSourceInput(
	candidate SubtitleFulfillmentCandidate,
	sourcePath string,
	streamInventory string,
) storage.MediaComponentSourceInput {
	releaseTitle := candidate.ReleaseTitle
	releaseName := candidate.ReleaseTitle
	releaseID := candidate.DownloadURL
	metadata := candidate.Metadata
	return storage.MediaComponentSourceInput{
		SourceRole:      "subtitle",
		SourceFilePath:  sourcePath,
		ReleaseTitle:    &releaseTitle,
		ReleaseName:     &releaseName,
		ReleaseID:       &releaseID,
		SourceMetadata:  &metadata,
		StreamInventory: streamInventory,
	}
}

func subtitleFulfillmentModes(item storage.MediaItem, target storage.MediaProfileSubtitleTarget) []string {
	source := strings.TrimSpace(target.Source)
	if source == "" {
		source = "any"
	}
	modes := []string{}
	if source != "embedded" && !externalSubtitleExists(item, target.LanguageID, firstMediaFilePath(item)) {
		modes = append(modes, "provider")
	}
	if source != "external" && !embeddedSubtitleExists(item, target) {
		modes = append(modes, "alternateRelease")
	}
	return modes
}

func bestSubtitleFulfillmentCandidate(
	item storage.MediaItem,
	releases []storage.ReleaseCandidateInput,
	need SubtitleFulfillmentNeed,
) (SubtitleFulfillmentCandidate, bool) {
	query := strings.ToLower(strings.Join(nonEmptyStrings([]string{need.LanguageID, firstString(need.Formats)}), " "))
	for _, release := range releases {
		title := strings.ToLower(release.Title)
		if query != "" && !containsAllTokens(title, query) {
			continue
		}
		return SubtitleFulfillmentCandidate{
			Need:         need,
			ReleaseTitle: release.Title,
			IndexerName:  release.IndexerName,
			DownloadURL:  release.DownloadURL,
			Metadata:     subtitleFulfillmentMetadata(release, need),
		}, true
	}
	return SubtitleFulfillmentCandidate{}, false
}

func subtitleFulfillmentQuery(item storage.MediaItem, need SubtitleFulfillmentNeed) string {
	return strings.Join(nonEmptyStrings([]string{
		decisions.SearchQueryForMediaItem(item),
		need.LanguageID,
		firstString(need.Formats),
	}), " ")
}

func embeddedSubtitleExists(item storage.MediaItem, target storage.MediaProfileSubtitleTarget) bool {
	for _, source := range item.ComponentSources {
		for _, artifact := range source.Artifacts {
			if artifact.StreamType == "subtitle" && languageMatches(artifact.Language, target.LanguageID) {
				return true
			}
		}
		for _, stream := range subtitleInventoryStreams(source.StreamInventory) {
			if subtitleStreamMatchesTarget(stream, target) {
				return true
			}
		}
	}
	return false
}

type subtitleInventoryStream struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Format   string `json:"format"`
	Codec    string `json:"codec"`
}

func subtitleInventoryStreams(payload string) []subtitleInventoryStream {
	var list []subtitleInventoryStream
	if err := json.Unmarshal([]byte(payload), &list); err == nil {
		return list
	}
	var wrapped struct {
		Streams []subtitleInventoryStream `json:"streams"`
	}
	if err := json.Unmarshal([]byte(payload), &wrapped); err != nil {
		return nil
	}
	return wrapped.Streams
}

func subtitleStreamMatchesTarget(stream subtitleInventoryStream, target storage.MediaProfileSubtitleTarget) bool {
	if stream.Type != "subtitle" || languageMatchKey(stream.Language) != languageMatchKey(target.LanguageID) {
		return false
	}
	formats := normalizedSubtitleFormats(target.Formats)
	if len(formats) == 0 {
		return true
	}
	return stringListHas(formats, firstNonEmpty(stream.Format, stream.Codec))
}

func subtitleFulfillmentMetadata(release storage.ReleaseCandidateInput, need SubtitleFulfillmentNeed) string {
	payload := map[string]any{
		"kind":     "subtitleFulfillment",
		"release":  release.Title,
		"indexer":  release.IndexerName,
		"language": need.LanguageID,
		"formats":  need.Formats,
		"mode":     need.Mode,
	}
	data, _ := json.Marshal(payload)
	return string(data)
}

func normalizedSubtitleFormats(values []string) []string {
	formats := []string{}
	for _, value := range values {
		if normalized := normalizeSubtitleFormat(value); normalized != "" {
			formats = append(formats, normalized)
		}
	}
	return formats
}

func firstMediaFilePath(item storage.MediaItem) string {
	if len(item.FilePaths) == 0 {
		return ""
	}
	return item.FilePaths[0]
}

func firstString(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func containsAllTokens(value string, tokens string) bool {
	for _, token := range strings.Fields(tokens) {
		if !strings.Contains(value, token) {
			return false
		}
	}
	return true
}
