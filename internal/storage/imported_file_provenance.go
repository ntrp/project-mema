package storage

import (
	"context"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

var importProvenanceExtensionPattern = regexp.MustCompile(`(?i)\.(mkv|mp4|avi|mov|wmv|ts|m2ts|iso)$`)

func recordImportedFileProvenance(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	sourcePath string,
	filePath string,
	importKind string,
) error {
	return recordImportedFileProvenanceWithOriginalName(ctx, q, mediaItemID, sourcePath, filePath, importKind, "")
}

func recordImportedFileProvenanceWithOriginalName(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	sourcePath string,
	filePath string,
	importKind string,
	originalFileName string,
) error {
	releaseTitle, releaseGroup := importedReleaseTitleAndGroup(filePath)
	if releaseGroup == "" {
		return nil
	}
	originalPath := strings.TrimSpace(sourcePath)
	transformation := map[string]any{
		"kind":         importKind,
		"originalPath": originalPath,
		"targetPath":   filePath,
	}
	if name := strings.TrimSpace(originalFileName); name != "" {
		transformation["originalFileName"] = name
	}
	_, err := upsertMediaComponentProvenance(ctx, q, MediaComponentProvenanceInput{
		MediaItemID:         mediaItemID,
		ComponentType:       "container",
		ComponentKey:        importedContainerComponentKey(filePath),
		ReleaseGroup:        releaseGroup,
		ReleaseName:         releaseTitle,
		SourceFilePath:      &filePath,
		TransformationChain: []map[string]any{transformation},
	})
	return err
}

func importedReleaseTitleAndGroup(filePath string) (string, string) {
	base := filepath.Base(strings.TrimSpace(filePath))
	title := importProvenanceExtensionPattern.ReplaceAllString(base, "")
	index := strings.LastIndex(title, "-")
	if index <= 0 || index >= len(title)-1 {
		return title, ""
	}
	group := strings.TrimSpace(title[index+1:])
	if strings.ContainsAny(group, " ._") || len(group) > 40 {
		return title, ""
	}
	return strings.TrimSpace(title[:index]), group
}

func importedContainerComponentKey(filePath string) string {
	return "imported:" + filepath.Clean(filePath)
}
