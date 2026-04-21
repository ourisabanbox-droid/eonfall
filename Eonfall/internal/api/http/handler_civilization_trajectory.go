package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"project-eonfall/internal/world"
)

func (h *Handler) GetCivilizationTrajectory(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	civilizationID, err := uuid.Parse(chi.URLParam(r, "civilizationID"))
	if err != nil {
		http.Error(w, "invalid civilizationID", http.StatusBadRequest)
		return
	}

	t, err := h.trajectoryRepo.GetByCivilizationID(r.Context(), civilizationID)
	if err != nil {
		http.Error(w, "trajectory not found", http.StatusNotFound)
		return
	}

	if t.WorldID != worldID {
		http.Error(w, "trajectory not found", http.StatusNotFound)
		return
	}

	profile := world.ComputeCivilizationProfile(t)

	resp := CivilizationTrajectoryResponse{
		CivilizationID: t.CivilizationID.String(),
		WorldID:        t.WorldID.String(),
		Scores: CivilizationTrajectoryScoresResponse{
			Resilience: t.ResilienceScore,
			Expansion:  t.ExpansionScore,
			Influence:  t.InfluenceScore,
			Science:    t.ScienceScore,
		},
		DominantAxis:  profile.DominantAxis,
		SecondaryAxis: profile.SecondaryAxis,
		ProfileLabel:  profile.ProfileLabel,
	}

	writeJSON(w, http.StatusOK, resp)
}
