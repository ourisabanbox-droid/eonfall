package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AcceptMissionRequest struct {
	MissionType    string `json:"mission_type"`
	TargetRegionID string `json:"target_region_id"`
}

func (h *Handler) AcceptCivilizationMission(w http.ResponseWriter, r *http.Request) {
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

	var req AcceptMissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if req.MissionType == "" {
		http.Error(w, "mission_type is required", http.StatusBadRequest)
		return
	}

	var targetRegionID *uuid.UUID
	if req.TargetRegionID != "" {
		parsed, err := uuid.Parse(req.TargetRegionID)
		if err != nil {
			http.Error(w, "invalid target_region_id", http.StatusBadRequest)
			return
		}
		targetRegionID = &parsed
	}

	if err := h.civilizationMissionRepo.Accept(
		r.Context(),
		worldID,
		civilizationID,
		req.MissionType,
		targetRegionID,
	); err != nil {
		http.Error(w, "failed to accept mission", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"world_id":         worldID.String(),
		"civilization_id":  civilizationID.String(),
		"mission_type":     req.MissionType,
		"target_region_id": req.TargetRegionID,
		"status":           "accepted",
	})
}
