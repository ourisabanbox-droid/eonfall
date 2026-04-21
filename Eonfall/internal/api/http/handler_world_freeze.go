package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) FreezeWorld(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	if err := h.worldRepo.SetFrozen(r.Context(), worldID, true); err != nil {
		http.Error(w, "failed to freeze world", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"world_id":  worldID.String(),
		"is_frozen": true,
	})
}

func (h *Handler) UnfreezeWorld(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	if err := h.worldRepo.SetFrozen(r.Context(), worldID, false); err != nil {
		http.Error(w, "failed to unfreeze world", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"world_id":  worldID.String(),
		"is_frozen": false,
	})
}
