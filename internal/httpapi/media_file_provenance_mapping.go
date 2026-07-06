package httpapi

import (
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func hydrateTrackProvenance(
	path string,
	tracks []MediaFileTrack,
	provenance []storage.MediaComponentProvenance,
) {
	if len(tracks) == 0 || len(provenance) == 0 {
		return
	}
	for index := range tracks {
		record, ok := matchingTrackProvenance(path, tracks, tracks[index], provenance)
		if ok {
			tracks[index].Provenance = mediaFileTrackProvenance(record)
		}
	}
}

func matchingTrackProvenance(
	path string,
	tracks []MediaFileTrack,
	track MediaFileTrack,
	provenance []storage.MediaComponentProvenance,
) (storage.MediaComponentProvenance, bool) {
	if track.Index != nil {
		if record, ok := provenanceByStreamID(track, provenance); ok {
			return record, true
		}
	}
	if record, ok := singleProvenanceForTrackType(tracks, track, provenance); ok {
		return record, true
	}
	if record, ok := containerProvenanceForPath(path, provenance); ok {
		return record, true
	}
	return storage.MediaComponentProvenance{}, false
}

func provenanceByStreamID(
	track MediaFileTrack,
	provenance []storage.MediaComponentProvenance,
) (storage.MediaComponentProvenance, bool) {
	for _, record := range provenance {
		if sameTrackComponent(record, track) && record.SourceStreamID != nil && *record.SourceStreamID == *track.Index {
			return record, true
		}
	}
	return storage.MediaComponentProvenance{}, false
}

func singleProvenanceForTrackType(
	tracks []MediaFileTrack,
	track MediaFileTrack,
	provenance []storage.MediaComponentProvenance,
) (storage.MediaComponentProvenance, bool) {
	if trackTypeCount(tracks, track.Type) != 1 {
		return storage.MediaComponentProvenance{}, false
	}
	var match storage.MediaComponentProvenance
	count := 0
	for _, record := range provenance {
		if sameTrackComponent(record, track) {
			match = record
			count++
		}
	}
	return match, count == 1
}

func trackTypeCount(tracks []MediaFileTrack, trackType MediaFileTrackType) int {
	count := 0
	for _, track := range tracks {
		if track.Type == trackType {
			count++
		}
	}
	return count
}

func sameTrackComponent(record storage.MediaComponentProvenance, track MediaFileTrack) bool {
	return strings.EqualFold(record.ComponentType, string(track.Type))
}

func containerProvenanceForPath(
	path string,
	provenance []storage.MediaComponentProvenance,
) (storage.MediaComponentProvenance, bool) {
	var match storage.MediaComponentProvenance
	count := 0
	for _, record := range provenance {
		if !strings.EqualFold(record.ComponentType, "container") || record.SourceFilePath == nil {
			continue
		}
		if *record.SourceFilePath == path {
			match = record
			count++
		}
	}
	return match, count == 1
}

func mediaFileTrackProvenance(record storage.MediaComponentProvenance) *MediaFileTrackProvenance {
	return &MediaFileTrackProvenance{
		Id:                  openapi_types.UUID(record.ID),
		MediaItemId:         openapi_types.UUID(record.MediaItemID),
		ComponentType:       record.ComponentType,
		ComponentKey:        record.ComponentKey,
		ReleaseGroup:        record.ReleaseGroup,
		ReleaseName:         record.ReleaseName,
		ReleaseId:           record.ReleaseID,
		SourceProvider:      record.SourceProvider,
		SourceFilePath:      record.SourceFilePath,
		RetainedSourceId:    optionalOpenAPIUUID(record.RetainedSourceID),
		SourceStreamId:      record.SourceStreamID,
		TransformationChain: record.TransformationChain,
		CreatedAt:           record.CreatedAt,
		UpdatedAt:           record.UpdatedAt,
	}
}
