package httpapi

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	users, err := s.settings.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "user_list_failed", "Could not list users")
		return
	}

	response := UserListResponse{Users: make([]ManagedUser, 0, len(users))}
	for _, user := range users {
		response.Users = append(response.Users, managedUserResponse(user))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body UserCreateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := userCreateInput(w, body)
	if !ok {
		return
	}

	user, err := s.settings.CreateUser(r.Context(), input)
	if err != nil {
		writeUserError(w, err, "Could not create user")
		return
	}
	writeJSON(w, http.StatusCreated, managedUserResponse(user))
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body UserUpdateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := userUpdateInput(w, body)
	if !ok {
		return
	}

	user, err := s.settings.UpdateUser(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeUserError(w, err, "Could not update user")
		return
	}
	writeJSON(w, http.StatusOK, managedUserResponse(user))
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request, id ResourceId) {
	session, ok := s.requireAdmin(w, r)
	if !ok {
		return
	}
	if uuid.UUID(session.user.Id) == uuid.UUID(id) {
		writeError(w, http.StatusBadRequest, "cannot_delete_current_user", "You cannot delete the current user")
		return
	}

	if err := s.settings.DeleteUser(r.Context(), uuid.UUID(id)); err != nil {
		writeUserError(w, err, "Could not delete user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeUserError(w http.ResponseWriter, err error, message string) {
	switch {
	case errors.Is(err, storage.ErrDuplicateUser):
		writeError(w, http.StatusBadRequest, "duplicate_user", "A user with that username already exists")
	case errors.Is(err, storage.ErrLastAdmin):
		writeError(w, http.StatusBadRequest, "last_admin_required", "At least one admin user is required")
	case errors.Is(err, storage.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", message)
	default:
		writeError(w, http.StatusInternalServerError, "user_update_failed", message)
	}
}
