package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) ListDLNARendererProfiles(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	profiles, err := s.settings.ListDLNARendererProfiles(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "dlna_profiles_failed", "Could not list DLNA profiles")
		return
	}
	writeJSON(w, http.StatusOK, dlnaRendererProfileListResponse(profiles))
}

func (s *Server) GetDLNARendererProfile(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	profile, err := s.settings.GetDLNARendererProfile(r.Context(), id)
	if err != nil {
		writeSettingsError(w, err, "Could not find DLNA profile")
		return
	}
	writeJSON(w, http.StatusOK, dlnaRendererProfileResponse(profile))
}

func (s *Server) CreateDLNARendererProfile(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body DLNARendererProfileCreateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, err := dlnaRendererProfileCreateInput(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_profile", "DLNA profile payload is invalid")
		return
	}
	profile, err := s.settings.CreateDLNARendererProfile(r.Context(), body.Id, input)
	if err != nil {
		writeSettingsError(w, err, "Could not create DLNA profile")
		return
	}
	s.refreshDLNARendererProfiles(r)
	writeJSON(w, http.StatusCreated, dlnaRendererProfileResponse(profile))
}

func (s *Server) UpdateDLNARendererProfile(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body DLNARendererProfileRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, err := dlnaRendererProfileInput(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_profile", "DLNA profile payload is invalid")
		return
	}
	profile, err := s.settings.UpdateDLNARendererProfile(r.Context(), id, input)
	if err != nil {
		writeSettingsError(w, err, "Could not update DLNA profile")
		return
	}
	s.refreshDLNARendererProfiles(r)
	writeJSON(w, http.StatusOK, dlnaRendererProfileResponse(profile))
}

func (s *Server) DeleteDLNARendererProfile(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	if err := s.settings.DeleteDLNARendererProfile(r.Context(), id); err != nil {
		writeSettingsError(w, err, "Could not delete DLNA profile")
		return
	}
	s.refreshDLNARendererProfiles(r)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ResetDLNARendererProfile(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	profile, err := s.settings.ResetDLNARendererProfile(r.Context(), id)
	if err != nil {
		writeSettingsError(w, err, "Could not reset DLNA profile")
		return
	}
	s.refreshDLNARendererProfiles(r)
	writeJSON(w, http.StatusOK, dlnaRendererProfileResponse(profile))
}

func (s *Server) CloneDLNARendererProfile(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body DLNARendererProfileCloneRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	profile, err := s.settings.CloneDLNARendererProfile(r.Context(), id, body.Id, body.Name)
	if err != nil {
		writeSettingsError(w, err, "Could not clone DLNA profile")
		return
	}
	s.refreshDLNARendererProfiles(r)
	writeJSON(w, http.StatusCreated, dlnaRendererProfileResponse(profile))
}

func (s *Server) ExportDLNARendererProfile(w http.ResponseWriter, r *http.Request, id string) {
	s.GetDLNARendererProfile(w, r, id)
}

func (s *Server) ImportDLNARendererProfile(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body DLNARendererProfileCreateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, err := dlnaRendererProfileCreateInput(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_profile", "DLNA profile payload is invalid")
		return
	}
	profile, err := s.settings.ImportDLNARendererProfile(r.Context(), body.Id, input)
	if err != nil {
		writeSettingsError(w, err, "Could not import DLNA profile")
		return
	}
	s.refreshDLNARendererProfiles(r)
	writeJSON(w, http.StatusOK, dlnaRendererProfileResponse(profile))
}

func (s *Server) ListDLNARendererDeviceOverrides(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	overrides, err := s.settings.ListDLNARendererDeviceOverrides(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "dlna_overrides_failed", "Could not list DLNA overrides")
		return
	}
	writeJSON(w, http.StatusOK, dlnaRendererOverrideListResponse(overrides))
}

func (s *Server) UpsertDLNARendererDeviceOverride(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body DLNARendererDeviceOverrideRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, err := dlnaRendererOverrideInput(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_override", "DLNA override payload is invalid")
		return
	}
	override, err := s.settings.UpsertDLNARendererDeviceOverride(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not save DLNA override")
		return
	}
	s.refreshDLNARendererProfiles(r)
	writeJSON(w, http.StatusOK, dlnaRendererOverrideResponse(override))
}

func (s *Server) DeleteDLNARendererDeviceOverride(
	w http.ResponseWriter,
	r *http.Request,
	id openapi_types.UUID,
) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	if err := s.settings.DeleteDLNARendererDeviceOverride(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete DLNA override")
		return
	}
	s.refreshDLNARendererProfiles(r)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ListDLNARecentDevices(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	writeJSON(w, http.StatusOK, DLNARecentDeviceListResponse{
		Devices: dlnaClientDiagnostics(s.currentDLNAStatus().RecentClients),
	})
}

func (s *Server) refreshDLNARendererProfiles(r *http.Request) {
	if s.dlna != nil {
		_ = s.dlna.RefreshRendererProfiles(r.Context())
	}
}
