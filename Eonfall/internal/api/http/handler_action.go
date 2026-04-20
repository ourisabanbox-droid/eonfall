package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetWorldActions(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v <= 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = v
	}

	actions, err := h.actionRepo.ListByWorldID(r.Context(), worldID, limit)
	if err != nil {
		http.Error(w, "failed to load actions", http.StatusInternalServerError)
		return
	}

	resp := make([]ActionResponse, 0, len(actions))
	for _, action := range actions {
		resp = append(resp, toActionResponse(action))
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"actions": resp,
	})
}
