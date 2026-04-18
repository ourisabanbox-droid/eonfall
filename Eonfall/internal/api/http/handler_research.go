package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"project-eonfall/internal/world"
)

func (h *Handler) PostResearchAction(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	var req ResearchActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if err := ValidateResearchActionRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentTick, err := h.worldRepo.GetCurrentTick(r.Context(), worldID)
	if err != nil {
		http.Error(w, "world not found", http.StatusNotFound)
		return
	}

	action := &world.WorldAction{
		ID:             uuid.New(),
		WorldID:        worldID,
		CivilizationID: req.CivilizationID,
		UserID:         req.UserID,
		ActionType:     world.ActionResearchStart,
		State:          world.ActionPending,
		TargetTick:     currentTick + 1,
		Payload: map[string]any{
			"research_type": req.ResearchType,
		},
		CreatedAt: time.Now().UTC(),
	}

	if err := h.actionRepo.Enqueue(r.Context(), action); err != nil {
		http.Error(w, "failed to enqueue action", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"action_id":   action.ID,
		"world_id":    action.WorldID,
		"target_tick": action.TargetTick,
		"action_type": action.ActionType,
		"status":      action.State,
	})
}
