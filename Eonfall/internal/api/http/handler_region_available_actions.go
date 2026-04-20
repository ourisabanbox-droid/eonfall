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

	actions := make([]AvailableActionResponse, 0, 2)

	if region.OwnerCivilizationID != nil {
		stabilizeReasons := []string{}
		stabilizeRecommended := false

		if region.RevoltRisk >= 60 {
			stabilizeRecommended = true
			stabilizeReasons = append(stabilizeReasons, "revolt_risk_high")
		}
		if region.Stability <= 30 {
			stabilizeRecommended = true
			stabilizeReasons = append(stabilizeReasons, "stability_low")
		}

		actions = append(actions, AvailableActionResponse{
			ActionType:  "stabilize_region",
			Label:       "Stabiliser la région",
			Description: "Réduit le risque de révolte et améliore temporairement la stabilité.",
			Recommended: stabilizeRecommended,
			Reasons:     stabilizeReasons,
		})

		droughtReasons := []string{}
		droughtRecommended := false

		if region.DroughtRisk >= 60 {
			droughtRecommended = true
			droughtReasons = append(droughtReasons, "drought_risk_high")
		}

		actions = append(actions, AvailableActionResponse{
			ActionType:  "drought_relief",
			Label:       "Secours hydrique",
			Description: "Réduit le risque de sécheresse et peut réinjecter du stock alimentaire.",
			Recommended: droughtRecommended,
			Reasons:     droughtReasons,
		})
	}

	signals := AvailableActionSignals{
		Stability:   region.Stability,
		RevoltRisk:  region.RevoltRisk,
		DroughtRisk: region.DroughtRisk,
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"actions":         actions,
		"current_signals": signals,
	})
}
