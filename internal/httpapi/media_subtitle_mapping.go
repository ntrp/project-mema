package httpapi

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func mediaSubtitleResponses(values []storage.MediaItemSubtitle) []MediaItemSubtitle {
	items := make([]MediaItemSubtitle, 0, len(values))
	for _, value := range values {
		items = append(items, MediaItemSubtitle{
			Id:                 openapi_types.UUID(value.ID),
			SeasonId:           optionalOpenAPIUUID(value.SeasonID),
			EpisodeId:          optionalOpenAPIUUID(value.EpisodeID),
			ProviderId:         optionalOpenAPIUUID(value.ProviderID),
			ProviderName:       value.ProviderName,
			LanguageId:         value.LanguageID,
			Format:             value.Format,
			FilePath:           value.FilePath,
			SourceUrl:          value.SourceURL,
			SourceReference:    value.SourceRef,
			ReleaseName:        value.ReleaseName,
			ProviderSubtitleId: value.ProviderSubtitleID,
			Checksum:           value.Checksum,
			SizeBytes:          value.SizeBytes,
			DownloadedAt:       value.DownloadedAt,
		})
	}
	return items
}
