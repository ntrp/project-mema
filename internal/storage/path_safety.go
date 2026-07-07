package storage

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func safeAbsRoot(root string) (string, error) {
	root = filepath.Clean(strings.TrimSpace(root))
	if root == "" || root == "." || !filepath.IsAbs(root) || root == string(os.PathSeparator) {
		return "", ErrInvalidInput
	}
	return root, nil
}

func absoluteCleanPath(path string) (string, error) {
	path = filepath.Clean(strings.TrimSpace(path))
	if path == "" || path == "." {
		return "", ErrInvalidInput
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", ErrInvalidInput
	}
	return absolute, nil
}

func absoluteCleanPathOrClean(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	absolute, err := absoluteCleanPath(path)
	if err == nil {
		return absolute
	}
	return filepath.Clean(strings.TrimSpace(path))
}

func safePathUnderRoot(root string, value string, allowRoot bool) (string, error) {
	root, err := safeAbsRoot(root)
	if err != nil {
		return "", err
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidInput
	}
	if encodedTraversal(value) {
		return "", ErrInvalidInput
	}
	target := filepath.Clean(value)
	if !filepath.IsAbs(target) {
		target = filepath.Join(root, target)
	}
	if err := validatePathInRoot(root, target, allowRoot); err != nil {
		return "", err
	}
	return target, nil
}

func validatePathInRoot(root string, target string, allowRoot bool) error {
	root, err := safeAbsRoot(root)
	if err != nil {
		return err
	}
	target = filepath.Clean(strings.TrimSpace(target))
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return ErrInvalidInput
	}
	if rel == "." {
		if allowRoot {
			return nil
		}
		return ErrInvalidInput
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return ErrInvalidInput
	}
	return nil
}

func encodedTraversal(value string) bool {
	decoded, err := url.PathUnescape(value)
	if err != nil || decoded == value {
		return err != nil
	}
	if filepath.IsAbs(decoded) {
		return true
	}
	for _, part := range strings.Split(filepath.ToSlash(filepath.Clean(decoded)), "/") {
		if part == ".." {
			return true
		}
	}
	return false
}
