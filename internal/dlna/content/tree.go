package content

import (
	"context"
	"errors"
	"sort"
	"strconv"

	"media-manager/internal/storage"
)

var ErrObjectNotFound = errors.New("dlna content object not found")

func (t *Tree) BrowseChildren(ctx context.Context, id string) ([]Object, error) {
	items, err := t.visibleItems(ctx)
	if err != nil {
		return nil, err
	}
	ref, err := DecodeID(id)
	if err != nil {
		return nil, ErrObjectNotFound
	}
	switch {
	case id == RootID:
		return t.rootChildren(items), nil
	case ref.Kind == "root":
		return t.rootContainerChildren(ref.Key, items), nil
	case ref.Kind == "collection":
		return t.groupedMediaChildren(id, items, matchCollection(ref.Key)), nil
	case ref.Kind == "genre":
		return t.groupedMediaChildren(id, items, matchGenre(ref.Key)), nil
	case ref.Kind == "year":
		return t.groupedMediaChildren(id, items, matchYear(ref.Key)), nil
	case ref.Kind == "media":
		return t.mediaChildren(id, ref.Key, items), nil
	case ref.Kind == "season":
		return t.seasonChildren(id, ref.Key, items), nil
	case ref.Kind == "episode":
		return t.episodeChildren(id, ref.Key, items), nil
	case ref.Kind == "file":
		return []Object{}, nil
	default:
		return nil, ErrObjectNotFound
	}
}

func (t *Tree) BrowseMetadata(ctx context.Context, id string) (Object, error) {
	if id == RootID {
		return Object{ID: RootID, Title: "Mema", Class: "object.container", Kind: ObjectContainer}, nil
	}
	children, err := t.allObjects(ctx)
	if err != nil {
		return Object{}, err
	}
	for _, object := range children {
		if object.ID == id {
			return object, nil
		}
	}
	return Object{}, ErrObjectNotFound
}

func (t *Tree) rootChildren(items []storage.MediaItem) []Object {
	candidates := []Object{
		rootContainer("movies", "Movies", countItems(items, "movie")),
		rootContainer("series", "TV Shows", countItems(items, "serie")),
		rootContainer("collections", "Collections", len(collections(items))),
		rootContainer("recently-added", "Recently Added", len(items)),
		rootContainer("recently-updated", "Recently Updated", len(items)),
		rootContainer("genres", "Genres", len(genres(items))),
		rootContainer("years", "Years", len(years(items))),
	}
	children := make([]Object, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.ChildCount > 0 {
			children = append(children, candidate)
		}
	}
	return children
}

func (t *Tree) rootContainerChildren(key string, items []storage.MediaItem) []Object {
	id := EncodeID(RootContainerRef(key))
	switch key {
	case "movies":
		return t.mediaItemObjects(id, filterItems(items, "movie"))
	case "series":
		return t.mediaItemObjects(id, filterItems(items, "serie"))
	case "collections":
		return collectionChildren(id, items)
	case "recently-added":
		return t.recentChildren(id, items, true)
	case "recently-updated":
		return t.recentChildren(id, items, false)
	case "genres":
		return genreChildren(id, items)
	case "years":
		return yearChildren(id, items)
	default:
		return nil
	}
}

func (t *Tree) visibleItems(ctx context.Context) ([]storage.MediaItem, error) {
	items, err := t.source.ListMediaItems(ctx)
	if err != nil {
		return nil, err
	}
	visible := make([]storage.MediaItem, 0, len(items))
	for _, item := range items {
		if len(t.availableFiles(item)) > 0 {
			visible = append(visible, item)
		}
	}
	sortMediaItems(visible)
	return visible, nil
}

func (t *Tree) allObjects(ctx context.Context) ([]Object, error) {
	seen := map[string]struct{}{}
	var objects []Object
	var visit func(string) error
	visit = func(parentID string) error {
		children, err := t.BrowseChildren(ctx, parentID)
		if err != nil {
			return err
		}
		for _, child := range children {
			if _, ok := seen[child.ID]; ok {
				continue
			}
			seen[child.ID] = struct{}{}
			objects = append(objects, child)
			if err := visit(child.ID); err != nil {
				return err
			}
		}
		return nil
	}
	return objects, visit(RootID)
}

func (t *Tree) mediaItemObjects(parentID string, items []storage.MediaItem) []Object {
	objects := make([]Object, 0, len(items))
	for _, item := range items {
		files := t.availableFiles(item)
		count := len(files)
		if item.Type == "serie" {
			count = visibleSeasonCount(files, item.Seasons)
		}
		objects = append(objects, mediaObject(parentID, item, files, count))
	}
	return objects
}

func (t *Tree) groupedMediaChildren(parentID string, items []storage.MediaItem, match func(storage.MediaItem) bool) []Object {
	matched := make([]storage.MediaItem, 0, len(items))
	for _, item := range items {
		if match(item) {
			matched = append(matched, item)
		}
	}
	return t.mediaItemObjects(parentID, matched)
}

func countItems(items []storage.MediaItem, mediaType string) int {
	return len(filterItems(items, mediaType))
}

func filterItems(items []storage.MediaItem, mediaType string) []storage.MediaItem {
	filtered := make([]storage.MediaItem, 0, len(items))
	for _, item := range items {
		if item.Type == mediaType {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func sortMediaItems(items []storage.MediaItem) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Title == items[j].Title {
			return yearValue(items[i]) < yearValue(items[j])
		}
		return items[i].Title < items[j].Title
	})
}

func yearValue(item storage.MediaItem) int32 {
	if item.Year == nil {
		return 0
	}
	return *item.Year
}

func matchYear(year string) func(storage.MediaItem) bool {
	return func(item storage.MediaItem) bool {
		return item.Year != nil && strconv.Itoa(int(*item.Year)) == year
	}
}
