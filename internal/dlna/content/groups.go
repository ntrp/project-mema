package content

import (
	"sort"
	"strconv"
	"strings"

	"media-manager/internal/storage"
)

func collectionChildren(parentID string, items []storage.MediaItem) []Object {
	values := collections(items)
	objects := make([]Object, 0, len(values))
	for _, value := range values {
		objects = append(objects, groupContainer(parentID, "collection", value, countMatches(items, matchCollection(value))))
	}
	return objects
}

func genreChildren(parentID string, items []storage.MediaItem) []Object {
	values := genres(items)
	objects := make([]Object, 0, len(values))
	for _, value := range values {
		objects = append(objects, groupContainer(parentID, "genre", value, countMatches(items, matchGenre(value))))
	}
	return objects
}

func yearChildren(parentID string, items []storage.MediaItem) []Object {
	values := years(items)
	objects := make([]Object, 0, len(values))
	for _, value := range values {
		objects = append(objects, groupContainer(parentID, "year", value, countMatches(items, matchYear(value))))
	}
	return objects
}

func collections(items []storage.MediaItem) []string {
	set := map[string]struct{}{}
	for _, item := range items {
		if item.CollectionName != nil && strings.TrimSpace(*item.CollectionName) != "" {
			set[*item.CollectionName] = struct{}{}
		}
	}
	return sortedKeys(set)
}

func genres(items []storage.MediaItem) []string {
	set := map[string]struct{}{}
	for _, item := range items {
		for _, genre := range item.Genres {
			if strings.TrimSpace(genre) != "" {
				set[genre] = struct{}{}
			}
		}
	}
	return sortedKeys(set)
}

func years(items []storage.MediaItem) []string {
	set := map[string]struct{}{}
	for _, item := range items {
		if item.Year != nil {
			set[strconv.Itoa(int(*item.Year))] = struct{}{}
		}
	}
	return sortedKeys(set)
}

func matchCollection(value string) func(storage.MediaItem) bool {
	return func(item storage.MediaItem) bool {
		return item.CollectionName != nil && *item.CollectionName == value
	}
}

func matchGenre(value string) func(storage.MediaItem) bool {
	return func(item storage.MediaItem) bool {
		for _, genre := range item.Genres {
			if genre == value {
				return true
			}
		}
		return false
	}
}

func countMatches(items []storage.MediaItem, match func(storage.MediaItem) bool) int {
	count := 0
	for _, item := range items {
		if match(item) {
			count++
		}
	}
	return count
}

func sortedKeys(set map[string]struct{}) []string {
	values := make([]string, 0, len(set))
	for value := range set {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}
