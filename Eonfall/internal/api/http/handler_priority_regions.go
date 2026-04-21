package http

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetPriorityRegions(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	limit := 10
	if raw := r.URL.Query().Get("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v <= 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = v
	}

	recommendedOnly := false
	if raw := r.URL.Query().Get("recommended_only"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			http.Error(w, "invalid recommended_only", http.StatusBadRequest)
			return
		}
		recommendedOnly = v
	}

	regions, err := h.regionRepo.ListByWorldID(r.Context(), worldID)
	if err != nil {
		http.Error(w, "failed to load regions", http.StatusInternalServerError)
		return
	}

	out := make([]PriorityRegionResponse, 0, len(regions))

	for _, region := range regions {
		score := region.RevoltRisk + region.DroughtRisk + maxIntHTTP(0, 100-region.Stability)

		var recommended []AvailableActionResponse
		var civilizationContext *CivilizationContextResponse
		actionable := false
		var blockers []string

		if region.OwnerCivilizationID == nil {
			actionable = false
			blockers = []string{"region_unowned"}
		} else {
			actionable = true
			recommended, civilizationContext = h.buildRegionAvailableActions(r.Context(), worldID, region)
		}

		if recommendedOnly {
			filtered := make([]AvailableActionResponse, 0, len(recommended))
			for _, action := range recommended {
				if action.Recommended {
					filtered = append(filtered, action)
				}
			}
			recommended = filtered
		}

		var topRecommendedAction *string

		hasStabilize := false
		hasDroughtRelief := false

		for _, action := range recommended {
			if !action.Recommended {
				continue
			}
			switch action.ActionType {
			case "stabilize_region":
				hasStabilize = true
			case "drought_relief":
				hasDroughtRelief = true
			}
		}

		if hasStabilize && (!hasDroughtRelief || region.RevoltRisk >= region.DroughtRisk) {
			v := "stabilize_region"
			topRecommendedAction = &v
		} else if hasDroughtRelief {
			v := "drought_relief"
			topRecommendedAction = &v
		}

		out = append(out, PriorityRegionResponse{
			RegionID:             region.ID.String(),
			Stability:            region.Stability,
			RevoltRisk:           region.RevoltRisk,
			DroughtRisk:          region.DroughtRisk,
			PriorityScore:        score,
			RecommendedActions:   recommended,
			Actionable:           actionable,
			ActionBlockers:       blockers,
			TopRecommendedAction: topRecommendedAction,
			CivilizationContext:  civilizationContext,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].PriorityScore > out[j].PriorityScore
	})

	if limit < len(out) {
		out = out[:limit]
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"regions": out,
	})
}

func maxIntHTTP(a, b int) int {
	if a > b {
		return a
	}
	return b
}
