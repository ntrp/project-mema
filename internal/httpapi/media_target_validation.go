package httpapi

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func boolDefault(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func (s *Server) validateMediaTarget(ctx context.Context, qualityProfileID *string, libraryFolderID *uuid.UUID) error {
	if qualityProfileID == nil || strings.TrimSpace(*qualityProfileID) == "" {
		return errMissingQualityProfile
	}
	exists, err := s.settings.MediaProfileExists(ctx, strings.TrimSpace(*qualityProfileID))
	if err != nil {
		return err
	}
	if !exists {
		return errUnsupportedQualityProfile
	}
	if libraryFolderID == nil {
		return errMissingLibraryFolder
	}
	exists, err = s.settings.LibraryFolderExists(ctx, *libraryFolderID)
	if err != nil {
		return err
	}
	if !exists {
		return storage.ErrNotFound
	}
	return nil
}

func writeMediaTargetError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errMissingQualityProfile):
		writeError(w, http.StatusBadRequest, "quality_profile_required", "Quality profile is required")
	case errors.Is(err, errUnsupportedQualityProfile):
		writeError(w, http.StatusBadRequest, "quality_profile_invalid", "Quality profile is not supported")
	case errors.Is(err, errMissingLibraryFolder):
		writeError(w, http.StatusBadRequest, "library_folder_required", "Library folder is required")
	case errors.Is(err, storage.ErrNotFound):
		writeError(w, http.StatusNotFound, "library_folder_not_found", "Library folder was not found")
	default:
		writeError(w, http.StatusInternalServerError, "media_target_validation_failed", "Could not validate media target")
	}
}

func writeMediaRequestError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrRequestClosed):
		writeError(w, http.StatusBadRequest, "media_request_closed", "Media request is no longer pending")
	case errors.Is(err, storage.ErrNotFound):
		writeError(w, http.StatusNotFound, "media_request_not_found", "Could not find media request")
	default:
		writeError(w, http.StatusInternalServerError, "media_request_update_failed", "Could not update media request")
	}
}

var (
	errMissingQualityProfile     = errors.New("quality profile is required")
	errUnsupportedQualityProfile = errors.New("quality profile is not supported")
	errMissingLibraryFolder      = errors.New("library folder is required")
)
