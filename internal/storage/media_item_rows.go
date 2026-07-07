package storage

import (
	"path/filepath"
	"strings"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func mediaItemRecordParams(id uuid.UUID, input MediaItemInput, payloads mediaMetadataPayloads, mediaFolderPath *string) storagegen.CreateMediaItemRecordParams {
	return storagegen.CreateMediaItemRecordParams{
		ID:                  id,
		MediaType:           input.Type,
		ContentKind:         input.ContentKind,
		Title:               input.Title,
		Year:                int4Value(input.Year),
		Monitored:           input.Monitored,
		ExternalProvider:    textValue(input.ExternalProvider),
		ExternalID:          textValue(input.ExternalID),
		Overview:            textValue(input.Overview),
		PosterPath:          textValue(input.PosterPath),
		CollectionID:        textValue(input.CollectionID),
		CollectionName:      textValue(input.CollectionName),
		BackdropPath:        textValue(input.BackdropPath),
		MetadataStatus:      textValue(input.MetadataStatus),
		OriginalLanguage:    textValue(input.OriginalLanguage),
		SeriesType:          textValue(input.SeriesType),
		NumberingStrategy:   textValue(input.NumberingStrategy),
		ReleaseDate:         textValue(input.ReleaseDate),
		FirstAirDate:        textValue(input.FirstAirDate),
		RuntimeMinutes:      int4Value(input.RuntimeMinutes),
		SeasonCount:         int4Value(input.SeasonCount),
		EpisodeCount:        int4Value(input.EpisodeCount),
		VoteAverage:         float8Value(input.VoteAverage),
		Genres:              payloads.genres,
		Keywords:            payloads.keywords,
		Facts:               payloads.facts,
		Seasons:             payloads.seasons,
		CastMembers:         payloads.cast,
		CrewMembers:         payloads.crew,
		Recommendations:     payloads.recommendations,
		SimilarMedia:        payloads.similar,
		MonitorMode:         input.MonitorMode,
		MinimumAvailability: input.MinimumAvailability,
		QualityProfileID:    textValue(input.QualityProfileID),
		LibraryFolderID:     input.LibraryFolderID,
		MediaFolderPath:     textValue(mediaFolderPath),
	}
}

func mediaItemFromListRow(row storagegen.ListMediaItemsRow) MediaItem {
	return mediaItemFromGetRow(storagegen.GetMediaItemRow(row))
}

func mediaItemFromSearchRow(row storagegen.SearchMediaItemsRow) MediaItem {
	return mediaItemFromGetRow(storagegen.GetMediaItemRow(row))
}

func mediaItemFromMissingRow(row storagegen.ListMissingMediaItemsRow) MediaItem {
	return mediaItemFromGetRow(storagegen.GetMediaItemRow(row))
}

func mediaItemFromMatchRow(row storagegen.FindMonitoredMediaMatchRow) MediaItem {
	return mediaItemFromGetRow(storagegen.GetMediaItemRow(row))
}

func mediaItemFromGetRow(row storagegen.GetMediaItemRow) MediaItem {
	item := MediaItem{
		ID:                  row.ID,
		Type:                row.MediaType,
		ContentKind:         "standard",
		Title:               row.Title,
		Year:                int4Ptr(row.Year),
		Monitored:           row.Monitored,
		ExternalProvider:    textPtr(row.ExternalProvider),
		ExternalID:          textPtr(row.ExternalID),
		Overview:            textPtr(row.Overview),
		PosterPath:          textPtr(row.PosterPath),
		MonitorMode:         row.MonitorMode,
		SeriesType:          textPtr(row.SeriesType),
		MinimumAvailability: row.MinimumAvailability,
		QualityProfileID:    textPtr(row.QualityProfileID),
		QualityProfileName:  textPtr(row.QualityProfileName),
		Status:              row.Status,
		LibraryFolderID:     row.LibraryFolderID,
		MediaFolderPath:     effectiveMediaFolderPath(textPtr(row.MediaFolderPath), row.FilePaths),
		LibraryFolderPath:   emptyStringPtr(absoluteCleanPathOrClean(row.LibraryFolderPath)),
		FilePaths:           absoluteCleanPaths(row.FilePaths),
		Tags:                row.Tags,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
	item.CollectionID = textPtr(row.CollectionID)
	item.CollectionName = textPtr(row.CollectionName)
	item.BackdropPath = textPtr(row.BackdropPath)
	item.MetadataStatus = textPtr(row.MetadataStatus)
	item.OriginalLanguage = textPtr(row.OriginalLanguage)
	item.ReleaseDate = textPtr(row.ReleaseDate)
	item.FirstAirDate = textPtr(row.FirstAirDate)
	item.RuntimeMinutes = int4Ptr(row.RuntimeMinutes)
	item.SeasonCount = int4Ptr(row.SeasonCount)
	item.EpisodeCount = int4Ptr(row.EpisodeCount)
	item.VoteAverage = float8Ptr(row.VoteAverage)
	scanMediaMetadata(&item.MediaMetadataSnapshot, row.Genres, row.Keywords, row.Facts, row.Seasons, row.CastMembers, row.CrewMembers, row.Recommendations, row.SimilarMedia)
	item.MetadataFilePaths = collectMetadataFilePaths(item.FilePaths)
	return item
}

func effectiveMediaFolderPath(stored *string, filePaths []string) *string {
	stored = absoluteStringPtr(stored)
	filePaths = absoluteCleanPaths(filePaths)
	if stored != nil && mediaRootContainsFiles(*stored, filePaths) {
		return stored
	}
	currentRoot, ok := commonMediaFileRoot(filePaths)
	if ok {
		return &currentRoot
	}
	return stored
}

func absoluteStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	absolute := absoluteCleanPathOrClean(*value)
	if strings.TrimSpace(absolute) == "" {
		return nil
	}
	return &absolute
}

func absoluteCleanPaths(paths []string) []string {
	cleaned := make([]string, 0, len(paths))
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		cleaned = append(cleaned, absoluteCleanPathOrClean(path))
	}
	return cleaned
}

func mediaRootContainsFiles(root string, filePaths []string) bool {
	root = filepath.Clean(strings.TrimSpace(root))
	if root == "" || root == "." {
		return false
	}
	hasFile := false
	for _, path := range filePaths {
		path = filepath.Clean(strings.TrimSpace(path))
		if path == "" || path == "." {
			continue
		}
		hasFile = true
		if !pathWithinOrEqual(root, path) || root == path {
			return false
		}
	}
	return hasFile
}

func commonMediaFileRoot(filePaths []string) (string, bool) {
	root := ""
	for _, path := range filePaths {
		path = filepath.Clean(strings.TrimSpace(path))
		if path == "" || path == "." {
			continue
		}
		dir := filepath.Dir(path)
		if dir == "." || dir == string(filepath.Separator) {
			continue
		}
		if root == "" {
			root = dir
			continue
		}
		for root != "" && root != string(filepath.Separator) && !pathWithinOrEqual(root, dir) {
			root = filepath.Dir(root)
		}
	}
	if root == "" || root == "." || root == string(filepath.Separator) {
		return "", false
	}
	return root, true
}

func pathWithinOrEqual(root string, target string) bool {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func emptyStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
