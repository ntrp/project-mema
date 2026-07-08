package httpapi

import (
	"net/http"

	"media-manager/internal/targets"
)

func (s *Server) ListManualFulfillmentActions(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	actions := targets.ManualActions()
	response := ManualFulfillmentActionsResponse{
		Actions: make([]ManualFulfillmentAction, 0, len(actions)),
	}
	for _, action := range actions {
		response.Actions = append(response.Actions, manualFulfillmentActionResponse(action))
	}
	writeJSON(w, http.StatusOK, response)
}

func manualFulfillmentActionResponse(action targets.ManualAction) ManualFulfillmentAction {
	return ManualFulfillmentAction{
		Id:            action.ID,
		Operation:     TargetOperationType(action.Operation),
		Label:         action.Label,
		Description:   action.Description,
		Manual:        action.Manual,
		Automatic:     action.Automatic,
		Available:     action.Available,
		BlockedReason: action.BlockedReason,
		Method:        action.Method,
		Path:          action.Path,
		WorkerPath:    action.WorkerPath,
		StateEffect:   action.StateEffect,
	}
}
