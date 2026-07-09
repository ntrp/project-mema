package httpapi

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/mediafacts"
	"media-manager/internal/storage"
)

func hydrateFileTrackIDs(item storage.MediaItem, path string, tracks []MediaFileTrack) {
	for index := range tracks {
		if tracks[index].Id != nil {
			continue
		}
		tracks[index].Id = mediaFileTrackID(item, path, tracks[index])
	}
}

func mediaFileTrackID(item storage.MediaItem, path string, track MediaFileTrack) *openapi_types.UUID {
	if id := persistedMediaFileTrackID(item.FileFacts, path, track); id != nil {
		return id
	}
	streamIndex := int32(0)
	if track.Index != nil {
		streamIndex = *track.Index
	}
	id := mediafacts.TrackID(item.ID, path, string(track.Type), streamIndex)
	value := openapi_types.UUID(id)
	return &value
}

func persistedMediaFileTrackID(
	facts []storage.MediaFileFact,
	path string,
	track MediaFileTrack,
) *openapi_types.UUID {
	if track.Index == nil {
		return nil
	}
	for _, fact := range facts {
		if fact.FilePath != path {
			continue
		}
		for _, stored := range fact.Tracks {
			if stored.StreamIndex != *track.Index || stored.TrackType != string(track.Type) {
				continue
			}
			value := openapi_types.UUID(stored.ID)
			return &value
		}
	}
	return nil
}
