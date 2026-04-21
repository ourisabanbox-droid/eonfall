package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"project-eonfall/internal/worldrepo"
)

type Handler struct {
	worldRepo        *worldrepo.WorldRepository
	civilizationRepo *worldrepo.CivilizationRepository
	regionRepo       *worldrepo.RegionRepository
	actionRepo       *worldrepo.ActionRepository
	researchRepo     *worldrepo.ResearchRepository
	alertRepo        *worldrepo.AlertRepository
	catastropheRepo  *worldrepo.CatastropheRepository
	trajectoryRepo   *worldrepo.CivilizationTrajectoryRepository
}

func NewHandler(
	worldRepo *worldrepo.WorldRepository,
	civilizationRepo *worldrepo.CivilizationRepository,
	regionRepo *worldrepo.RegionRepository,
	actionRepo *worldrepo.ActionRepository,
	researchRepo *worldrepo.ResearchRepository,
	alertRepo *worldrepo.AlertRepository,
	catastropheRepo *worldrepo.CatastropheRepository,
	trajectoryRepo *worldrepo.CivilizationTrajectoryRepository,
) *Handler {
	return &Handler{
		worldRepo:        worldRepo,
		civilizationRepo: civilizationRepo,
		regionRepo:       regionRepo,
		actionRepo:       actionRepo,
		researchRepo:     researchRepo,
		alertRepo:        alertRepo,
		catastropheRepo:  catastropheRepo,
		trajectoryRepo:   trajectoryRepo,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetWorldSummary(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	world, err := h.worldRepo.GetByID(r.Context(), worldID)
	if err != nil {
		http.Error(w, "world not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":            world.ID,
		"name":          world.Name,
		"state":         world.State,
		"phase":         world.Phase,
		"current_tick":  world.CurrentTick,
		"season_number": world.SeasonNumber,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) GetWorldMap(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	regions, err := h.regionRepo.ListByWorldID(r.Context(), worldID)
	if err != nil {
		http.Error(w, "cannot load regions", http.StatusInternalServerError)
		return
	}

	type regionView struct {
		ID                  uuid.UUID  `json:"id"`
		Q                   int        `json:"q"`
		R                   int        `json:"r"`
		Biome               string     `json:"biome"`
		OwnerCivilizationID *uuid.UUID `json:"owner_civilization_id"`
		Stability           int        `json:"stability"`
		DroughtRisk         int        `json:"drought_risk"`
		RevoltRisk          int        `json:"revolt_risk"`
		HasCapital          bool       `json:"has_capital"`
	}

	out := make([]regionView, 0, len(regions))
	for _, rg := range regions {
		out = append(out, regionView{
			ID:                  rg.ID,
			Q:                   rg.Q,
			R:                   rg.R,
			Biome:               rg.Biome,
			OwnerCivilizationID: rg.OwnerCivilizationID,
			Stability:           rg.Stability,
			DroughtRisk:         rg.DroughtRisk,
			RevoltRisk:          rg.RevoltRisk,
			HasCapital:          rg.HasCapital,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"world_id": worldID,
		"regions":  out,
	})
}

func (h *Handler) GetCivilizationDashboard(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	civID, err := uuid.Parse(chi.URLParam(r, "civilizationID"))
	if err != nil {
		http.Error(w, "invalid civilizationID", http.StatusBadRequest)
		return
	}

	civs, err := h.civilizationRepo.ListByWorldID(r.Context(), worldID)
	if err != nil {
		http.Error(w, "cannot load civilizations", http.StatusInternalServerError)
		return
	}

	for _, civ := range civs {
		if civ.ID == civID {
			writeJSON(w, http.StatusOK, map[string]any{
				"civilization_id": civ.ID,
				"name":            civ.Name,
				"template_id":     civ.TemplateID,
				"stocks": map[string]any{
					"food":      civ.FoodStock,
					"materials": civ.MaterialsStock,
					"energy":    civ.EnergyStock,
					"credit":    civ.CreditStock,
					"knowledge": civ.ScienceStock,
				},
				"cohesion":       civ.Cohesion,
				"influence":      civ.Influence,
				"military_score": civ.MilitaryScore,
				"victory_score":  civ.VictoryScore,
			})
			return
		}
	}

	http.Error(w, "civilization not found", http.StatusNotFound)
}
