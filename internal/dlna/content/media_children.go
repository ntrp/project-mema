package content

import (
	"sort"
	"strconv"

	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func (t *Tree) mediaChildren(parentID string, mediaID string, items []storage.MediaItem) []Object {
	for _, item := range items {
		if item.ID.String() != mediaID {
			continue
		}
		if item.Type == "serie" {
			return t.seriesChildren(parentID, item)
		}
		return t.fileChildrenForItem(parentID, item, t.availableFiles(item))
	}
	return nil
}

func (t *Tree) seriesChildren(parentID string, item storage.MediaItem) []Object {
	files := t.availableFiles(item)
	objects := make([]Object, 0, len(item.Seasons))
	for _, season := range item.Seasons {
		if season.ID == nil {
			continue
		}
		childCount := visibleEpisodeCount(files, season)
		if childCount == 0 {
			continue
		}
		objects = append(objects, seasonObject(parentID, season, childCount))
	}
	sort.SliceStable(objects, func(i, j int) bool {
		return objects[i].Title < objects[j].Title
	})
	return objects
}

func (t *Tree) seasonChildren(parentID string, seasonID string, items []storage.MediaItem) []Object {
	for _, item := range items {
		files := t.availableFiles(item)
		for _, season := range item.Seasons {
			if season.ID == nil || season.ID.String() != seasonID {
				continue
			}
			return t.episodeObjects(parentID, season, files)
		}
	}
	return nil
}

func (t *Tree) episodeChildren(parentID string, episodeID string, items []storage.MediaItem) []Object {
	for _, item := range items {
		files := t.availableFiles(item)
		for _, season := range item.Seasons {
			for _, episode := range season.Episodes {
				if episode.ID == nil || episode.ID.String() != episodeID {
					continue
				}
				matched := filesForEpisode(files, season.SeasonNumber, episode.EpisodeNumber)
				return t.fileChildrenForItem(parentID, item, matched)
			}
		}
	}
	return nil
}

func (t *Tree) episodeObjects(parentID string, season storage.MediaSeason, files []File) []Object {
	objects := make([]Object, 0, len(season.Episodes))
	for _, episode := range season.Episodes {
		if episode.ID == nil {
			continue
		}
		childCount := len(filesForEpisode(files, season.SeasonNumber, episode.EpisodeNumber))
		if childCount == 0 {
			continue
		}
		objects = append(objects, episodeObject(parentID, episode, childCount))
	}
	sort.SliceStable(objects, func(i, j int) bool {
		return episodeSortKey(objects[i]) < episodeSortKey(objects[j])
	})
	return objects
}

func (t *Tree) fileChildren(parentID string, mediaID uuid.UUID, files []File) []Object {
	objects := make([]Object, 0, len(files))
	for _, file := range files {
		objects = append(objects, fileObject(parentID, mediaID, file, nil))
	}
	sort.SliceStable(objects, func(i, j int) bool {
		return objects[i].Title < objects[j].Title
	})
	return objects
}

func (t *Tree) fileChildrenForItem(parentID string, item storage.MediaItem, files []File) []Object {
	objects := make([]Object, 0, len(files))
	for _, file := range files {
		objects = append(objects, fileObject(parentID, item.ID, file, SubtitlesForFile(item, file.Path)))
	}
	sort.SliceStable(objects, func(i, j int) bool {
		return objects[i].Title < objects[j].Title
	})
	return objects
}

func visibleSeasonCount(files []File, seasons []storage.MediaSeason) int {
	count := 0
	for _, season := range seasons {
		if visibleEpisodeCount(files, season) > 0 {
			count++
		}
	}
	return count
}

func visibleEpisodeCount(files []File, season storage.MediaSeason) int {
	count := 0
	for _, episode := range season.Episodes {
		if len(filesForEpisode(files, season.SeasonNumber, episode.EpisodeNumber)) > 0 {
			count++
		}
	}
	return count
}

func episodeSortKey(object Object) string {
	if object.EpisodeID == nil {
		return object.Title
	}
	return strconv.Itoa(len(object.Title)) + object.Title
}
