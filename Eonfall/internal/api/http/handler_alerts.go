package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetWorldAlerts(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}

	alerts, err := h.alertRepo.ListByWorldID(r.Context(), worldID, limit)
	if err != nil {
		http.Error(w, "failed to load alerts", http.StatusInternalServerError)
		return
	}

	type alertView struct {
		ID             uuid.UUID      `json:"id"`
		WorldID        uuid.UUID      `json:"world_id"`
		CivilizationID *uuid.UUID     `json:"civilization_id"`
		RegionID       *uuid.UUID     `json:"region_id"`
		AlertType      string         `json:"alert_type"`
		Severity       string         `json:"severity"`
		Title          string         `json:"title"`
		Message        string         `json:"message"`
		Payload        map[string]any `json:"payload"`
		CreatedAt      string         `json:"created_at"`
	}

	out := make([]alertView, 0, len(alerts))
	for _, a := range alerts {
		out = append(out, alertView{
			ID:             a.ID,
			WorldID:        a.WorldID,
			CivilizationID: a.CivilizationID,
			RegionID:       a.RegionID,
			AlertType:      a.AlertType,
			Severity:       a.Severity,
			Title:          a.Title,
			Message:        a.Message,
			Payload:        a.Payload,
			CreatedAt:      a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"world_id": worldID,
		"alerts":   out,
	})
}
