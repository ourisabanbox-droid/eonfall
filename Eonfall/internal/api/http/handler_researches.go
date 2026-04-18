package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetCivilizationResearches(w http.ResponseWriter, r *http.Request) {
	civilizationID, err := uuid.Parse(chi.URLParam(r, "civilizationID"))
	if err != nil {
		http.Error(w, "invalid civilizationID", http.StatusBadRequest)
		return
	}

	researches, err := h.researchRepo.ListByCivilizationID(r.Context(), civilizationID)
	if err != nil {
		http.Error(w, "failed to load researches", http.StatusInternalServerError)
		return
	}

	type researchView struct {
		ID            uuid.UUID `json:"id"`
		ResearchType  string    `json:"research_type"`
		State         string    `json:"state"`
		Progress      int       `json:"progress"`
		StartedTick   *int64    `json:"started_tick"`
		CompletedTick *int64    `json:"completed_tick"`
	}

	out := make([]researchView, 0, len(researches))
	for _, rp := range researches {
		out = append(out, researchView{
			ID:            rp.ID,
			ResearchType:  string(rp.ResearchType),
			State:         rp.State,
			Progress:      rp.Progress,
			StartedTick:   rp.StartedTick,
			CompletedTick: rp.CompletedTick,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"civilization_id": civilizationID,
		"researches":      out,
	})
}
