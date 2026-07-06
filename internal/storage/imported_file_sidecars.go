package storage

import (
	"context"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

func recordImportedFileSidecars(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	mediaPath string,
	seasonID *uuid.UUID,
	episodeID *uuid.UUID,
	subtitlePreferredMode string,
) error {
	for _, sidecar := range MediaSidecarsForFile(mediaPath) {
		if err := recordImportedFileSidecar(ctx, q, mediaItemID, mediaPath, seasonID, episodeID, subtitlePreferredMode, sidecar); err != nil {
			return err
		}
	}
	return nil
}

func recordImportedFileSidecar(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	mediaPath string,
	seasonID *uuid.UUID,
	episodeID *uuid.UUID,
	subtitlePreferredMode string,
	sidecar MediaSidecar,
) error {
	if _, err := upsertMediaItemSidecar(ctx, q, MediaItemSidecarInput{
		MediaItemID:   mediaItemID,
		MediaFilePath: mediaPath,
		FilePath:      sidecar.Path,
		SidecarType:   sidecar.Type,
		LanguageID:    sidecar.LanguageID,
		Format:        sidecar.Format,
	}); err != nil {
		return err
	}
	if sidecar.Type != MediaSidecarSubtitle || sidecar.LanguageID == "" {
		return nil
	}
	_, err := upsertMediaItemSubtitle(ctx, q, MediaItemSubtitleInput{
		MediaItemID:   mediaItemID,
		SeasonID:      seasonID,
		EpisodeID:     episodeID,
		ProviderName:  "library",
		LanguageID:    sidecar.LanguageID,
		Format:        sidecar.Format,
		FilePath:      sidecar.Path,
		ReleaseName:   sidecarString(filepath.Base(mediaPath)),
		DownloadedAt:  time.Now().UTC(),
		RetentionMode: importedSubtitleRetentionMode(subtitlePreferredMode),
	})
	return err
}

func importedSubtitleRetentionMode(subtitlePreferredMode string) SubtitleRetentionMode {
	if subtitlePreferredMode == "embedded" {
		return SubtitleRetentionMux
	}
	return SubtitleRetentionExternal
}

func sidecarString(value string) *string {
	return &value
}
