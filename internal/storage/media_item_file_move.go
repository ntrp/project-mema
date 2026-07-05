package storage

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func moveMediaItemFiles(item MediaItem, newRoot string) error {
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return nil
	}
	oldRoot, err := safeAbsRoot(*item.MediaFolderPath)
	if err != nil {
		return err
	}
	newRoot, err = safeAbsRoot(newRoot)
	if err != nil {
		return err
	}
	if oldRoot == newRoot {
		return nil
	}
	for _, source := range mediaItemMovePaths(item) {
		target, err := movedMediaFileTarget(oldRoot, newRoot, source)
		if err != nil {
			return err
		}
		if err := moveFile(source, target); err != nil {
			return err
		}
	}
	removeEmptyMediaDirs(oldRoot)
	return nil
}

func mediaItemMovePaths(item MediaItem) []string {
	seen := map[string]struct{}{}
	paths := []string{}
	for _, path := range append(item.FilePaths, item.MetadataFilePaths...) {
		cleaned := filepath.Clean(strings.TrimSpace(path))
		if cleaned == "." || cleaned == "" {
			continue
		}
		if _, ok := seen[cleaned]; ok {
			continue
		}
		seen[cleaned] = struct{}{}
		paths = append(paths, cleaned)
	}
	return paths
}

func movedMediaFileTarget(oldRoot string, newRoot string, source string) (string, error) {
	source, err := safePathUnderRoot(oldRoot, source, false)
	if err != nil {
		return "", err
	}
	oldRoot, err = safeAbsRoot(oldRoot)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(oldRoot, source)
	if err != nil {
		return "", ErrInvalidInput
	}
	return safePathUnderRoot(newRoot, rel, false)
}

func moveFile(source string, target string) error {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	if _, err := os.Stat(target); err == nil {
		return ErrInvalidInput
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	if err := os.Rename(source, target); err == nil {
		return nil
	}
	if err := copyFile(source, target); err != nil {
		return err
	}
	return os.Remove(source)
}

func copyFile(source string, target string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return err
	}
	output, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(output, input)
	return err
}

func removeEmptyMediaDirs(root string) {
	dirs := []string{}
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || !entry.IsDir() {
			return nil
		}
		dirs = append(dirs, path)
		return nil
	})
	for index := len(dirs) - 1; index >= 0; index-- {
		_ = os.Remove(dirs[index])
	}
}
