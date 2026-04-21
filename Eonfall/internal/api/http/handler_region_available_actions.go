package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetRegionAvailableActions(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	regionID, err := uuid.Parse(chi.URLParam(r, "regionID"))
	if err != nil {
		http.Error(w, "invalid regionID", http.StatusBadRequest)
		return
	}

	region, err := h.regionRepo.GetByID(r.Context(), worldID, regionID)
	if err != nil {
		http.Error(w, "region not found", http.StatusNotFound)
		return
	}

	actions, civilizationContext := h.buildRegionAvailableActions(r.Context(), worldID, region)

	signals := AvailableActionSignals{
		Stability:   region.Stability,
		RevoltRisk:  region.RevoltRisk,
		DroughtRisk: region.DroughtRisk,
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"actions":              actions,
		"current_signals":      signals,
		"civilization_context": civilizationContext,
	})
}
