package dlna

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	dlnaThumbnailCacheMaxAge   = 14 * 24 * time.Hour
	dlnaOutputCacheMaxAge      = 7 * 24 * time.Hour
	dlnaThumbnailCacheMaxBytes = 512 * 1024 * 1024
	dlnaOutputCacheMaxBytes    = 8 * 1024 * 1024 * 1024
)

type CachePruneResult struct {
	DeletedFiles int
	DeletedBytes int64
}

type cacheEntry struct {
	Path    string
	Size    int64
	ModTime time.Time
}

func (m *Manager) PruneCaches(now time.Time) (CachePruneResult, error) {
	thumbs, err := pruneCacheDir(m.thumbDir, now, dlnaThumbnailCacheMaxAge, dlnaThumbnailCacheMaxBytes)
	if err != nil {
		return thumbs, err
	}
	outputs, err := pruneCacheDir(m.remuxDir, now, dlnaOutputCacheMaxAge, dlnaOutputCacheMaxBytes)
	return CachePruneResult{
		DeletedFiles: thumbs.DeletedFiles + outputs.DeletedFiles,
		DeletedBytes: thumbs.DeletedBytes + outputs.DeletedBytes,
	}, err
}

func pruneCacheDir(dir string, now time.Time, maxAge time.Duration, maxBytes int64) (CachePruneResult, error) {
	entries, err := cacheEntries(dir)
	if err != nil {
		return CachePruneResult{}, err
	}
	result := CachePruneResult{}
	kept := make([]cacheEntry, 0, len(entries))
	total := int64(0)
	for _, entry := range entries {
		if maxAge > 0 && now.Sub(entry.ModTime) > maxAge {
			removeCacheEntry(entry, &result)
			continue
		}
		kept = append(kept, entry)
		total += entry.Size
	}
	sort.Slice(kept, func(i, j int) bool { return kept[i].ModTime.Before(kept[j].ModTime) })
	for _, entry := range kept {
		if maxBytes <= 0 || total <= maxBytes {
			break
		}
		if removeCacheEntry(entry, &result) {
			total -= entry.Size
		}
	}
	return result, nil
}

func cacheEntries(dir string) ([]cacheEntry, error) {
	entries := []cacheEntry{}
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		entries = append(entries, cacheEntry{Path: path, Size: info.Size(), ModTime: info.ModTime()})
		return nil
	})
	if os.IsNotExist(err) {
		return entries, nil
	}
	return entries, err
}

func removeCacheEntry(entry cacheEntry, result *CachePruneResult) bool {
	if os.Remove(entry.Path) != nil {
		return false
	}
	result.DeletedFiles++
	result.DeletedBytes += entry.Size
	return true
}
