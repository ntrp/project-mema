package dlna

import (
	"encoding/json"

	"media-manager/internal/dlna/content"
)

func DIDLOptionsForProfile(profile RendererProfile) content.DIDLOptions {
	options := profile.DIDLOptions
	if didlOptionsEmpty(options) {
		options = content.DefaultDIDLOptions()
	}
	if len(options.SubtitleFormats) == 0 {
		options.SubtitleFormats = append([]string{}, profile.SubtitleFormats...)
	}
	return options
}

func didlOptionsEmpty(options content.DIDLOptions) bool {
	return len(options.SubtitleFormats) == 0 &&
		!options.IncludeSubtitleResources &&
		!options.IncludeArtwork &&
		options.ArtworkProfileID == "" &&
		!options.IncludeDates &&
		!options.IncludeMediaMetadata &&
		!options.IncludeFolderData &&
		!options.IncludeChildCounts
}

func parseRendererArtwork(payload []byte) rendererArtworkJSON {
	var raw rendererArtworkJSON
	_ = json.Unmarshal(payload, &raw)
	return raw
}

func parseRendererMetadata(payload []byte) rendererMetadataJSON {
	var raw rendererMetadataJSON
	_ = json.Unmarshal(payload, &raw)
	return raw
}

func rendererDIDLOptions(
	subtitles rendererSubtitlesJSON,
	artwork rendererArtworkJSON,
	metadata rendererMetadataJSON,
) content.DIDLOptions {
	defaults := content.DefaultDIDLOptions()
	defaults.SubtitleFormats = append([]string{}, subtitles.Formats...)
	defaults.IncludeSubtitleResources = boolDefault(subtitles.Resources, true)
	defaults.IncludeArtwork = boolDefault(artwork.AlbumArt, boolDefault(metadata.AlbumArt, true))
	defaults.ArtworkProfileID = firstNonEmpty(artwork.ProfileID, "JPEG_TN")
	defaults.IncludeDates = boolDefault(metadata.Dates, true)
	defaults.IncludeMediaMetadata = boolDefault(metadata.Media, boolDefault(metadata.RichMetadata, true))
	defaults.IncludeFolderData = boolDefault(metadata.FolderData, true)
	defaults.IncludeChildCounts = boolDefault(metadata.ChildCounts, true)
	return defaults
}

func boolDefault(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}
