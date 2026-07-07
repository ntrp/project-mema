package httpapi

import (
	"net/http"

	"media-manager/internal/storage"
)

func (s *Server) GetDLNASettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	settings, err := s.settings.GetDLNASettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_load_failed", "Could not load DLNA settings")
		return
	}
	writeJSON(w, http.StatusOK, s.dlnaSettingsResponse(settings))
}

func (s *Server) UpdateDLNASettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body DLNASettingsRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	settings, err := s.settings.UpdateDLNASettings(r.Context(), storage.DLNASettingsInput{
		Enabled:                 body.Enabled,
		FriendlyName:            body.FriendlyName,
		Interfaces:              body.Interfaces,
		AllowedCIDRs:            body.AllowedCidrs,
		AnnounceIntervalSeconds: body.AnnounceIntervalSeconds,
		TranscodeEnabled:        body.TranscodeEnabled,
		ThumbnailsEnabled:       body.ThumbnailsEnabled,
		SubtitlesEnabled:        body.SubtitlesEnabled,
		DefaultRendererProfile:  body.DefaultRendererProfile,
	})
	if err != nil {
		writeSettingsError(w, err, "Could not update DLNA settings")
		return
	}
	if s.dlna != nil {
		if err := s.dlna.ApplySettings(r.Context(), settings); err != nil {
			writeError(w, http.StatusBadRequest, "dlna_start_failed", "Could not apply DLNA settings")
			return
		}
	}
	writeJSON(w, http.StatusOK, s.dlnaSettingsResponse(settings))
}

func (s *Server) RestartDLNA(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	settings, err := s.settings.GetDLNASettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_load_failed", "Could not load DLNA settings")
		return
	}
	if s.dlna != nil {
		if err := s.dlna.ApplySettings(r.Context(), settings); err != nil {
			writeError(w, http.StatusBadRequest, "dlna_start_failed", "Could not restart DLNA")
			return
		}
	}
	writeJSON(w, http.StatusOK, s.dlnaSettingsResponse(settings))
}
