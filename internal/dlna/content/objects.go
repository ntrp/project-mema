package content

import (
	"fmt"
	"path/filepath"

	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func rootContainer(key string, title string, childCount int) Object {
	id := EncodeID(RootContainerRef(key))
	return Object{
		ID:             id,
		ParentID:       RootID,
		Title:          title,
		Class:          "object.container",
		Kind:           ObjectContainer,
		ChildCount:     childCount,
		OmitChildCount: true,
	}
}

func groupContainer(parentID string, kind string, title string, childCount int) Object {
	return Object{
		ID:         EncodeID(GroupRef(kind, title)),
		ParentID:   parentID,
		Title:      title,
		Class:      "object.container.storageFolder",
		Kind:       ObjectContainer,
		ChildCount: childCount,
	}
}

func mediaObject(parentID string, item storage.MediaItem, files []File, childCount int) Object {
	id := item.ID
	class := "object.item.videoItem.movie"
	kind := ObjectItem
	if item.Type == "serie" {
		class = "object.container.album.videoAlbum"
		kind = ObjectContainer
	} else if len(files) > 1 {
		class = "object.container.storageFolder"
		kind = ObjectContainer
	}
	object := Object{
		ID:          EncodeID(MediaItemRef(item.ID)),
		ParentID:    parentID,
		Title:       mediaTitle(item),
		Class:       class,
		Kind:        kind,
		ChildCount:  childCount,
		MediaType:   item.Type,
		Year:        item.Year,
		Date:        mediaDate(item),
		Genres:      append([]string{}, item.Genres...),
		Artists:     mediaArtists(item),
		Album:       item.CollectionName,
		Artwork:     item.PosterPath,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		MediaItemID: &id,
	}
	if item.Type != "serie" && len(files) == 1 {
		object.ChildCount = 0
		object.FileHash = files[0].Hash
		object.FilePath = files[0].Path
		object.Subtitles = SubtitlesForFile(item, files[0].Path)
	}
	return object
}

func seasonObject(parentID string, season storage.MediaSeason, childCount int) Object {
	seasonID := uuid.Nil
	if season.ID != nil {
		seasonID = *season.ID
	}
	return Object{
		ID:         EncodeID(SeasonRef(seasonID)),
		ParentID:   parentID,
		Title:      seasonTitle(season),
		Class:      "object.container.album.videoAlbum",
		Kind:       ObjectContainer,
		ChildCount: childCount,
		SeasonID:   &seasonID,
	}
}

func episodeObject(parentID string, episode storage.MediaEpisode, childCount int) Object {
	episodeID := uuid.Nil
	if episode.ID != nil {
		episodeID = *episode.ID
	}
	return Object{
		ID:         EncodeID(EpisodeRef(episodeID)),
		ParentID:   parentID,
		Title:      episodeTitle(episode),
		Class:      "object.item.videoItem.episode",
		Kind:       ObjectItem,
		ChildCount: childCount,
		Date:       episode.AirDate,
		EpisodeID:  &episodeID,
	}
}

func fileObject(parentID string, mediaID uuid.UUID, file File, subtitles []Subtitle) Object {
	return Object{
		ID:          EncodeID(FileRef(mediaID, file.Path)),
		ParentID:    parentID,
		Title:       filepath.Base(file.Path),
		Class:       "object.item.videoItem",
		Kind:        ObjectItem,
		MediaItemID: &mediaID,
		FileHash:    file.Hash,
		FilePath:    file.Path,
		Subtitles:   subtitles,
	}
}

func mediaTitle(item storage.MediaItem) string {
	if item.Year == nil {
		return item.Title
	}
	return fmt.Sprintf("%s (%d)", item.Title, *item.Year)
}

func mediaDate(item storage.MediaItem) *string {
	if item.ReleaseDate != nil {
		return item.ReleaseDate
	}
	return item.FirstAirDate
}

func mediaArtists(item storage.MediaItem) []string {
	artists := make([]string, 0, len(item.Cast))
	for _, person := range item.Cast {
		if person.Name != "" {
			artists = append(artists, person.Name)
		}
	}
	return artists
}

func seasonTitle(season storage.MediaSeason) string {
	if season.Name != "" {
		return season.Name
	}
	return fmt.Sprintf("Season %d", season.SeasonNumber)
}

func episodeTitle(episode storage.MediaEpisode) string {
	if episode.Name != "" {
		return episode.Name
	}
	return fmt.Sprintf("Episode %d", episode.EpisodeNumber)
}
