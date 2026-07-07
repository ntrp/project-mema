package content

import (
	"sort"

	"media-manager/internal/storage"
)

func (t *Tree) recentChildren(parentID string, items []storage.MediaItem, added bool) []Object {
	sorted := append([]storage.MediaItem{}, items...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if added {
			return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
		}
		return sorted[i].UpdatedAt.After(sorted[j].UpdatedAt)
	})
	return t.mediaItemObjects(parentID, sorted)
}
